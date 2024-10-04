import Foundation
class FeedViewModel: ObservableObject {
    @Published var clips = [Clip]()
    @Published var isLoading = true
        
    init() {
        fetchClips { fetched
            in self.isLoading = !fetched
        }
    }

    func fetchClips(completion: @escaping (Bool) -> Void) {
        self.recommendVideos { fetchedClips in
            DispatchQueue.main.async {
                if !fetchedClips.isEmpty {
                    self.clips = fetchedClips
                    self.isLoading = false
                    completion(true) // Indicate success

                } else {
                    completion(false) // Indicate failure or decision not to load
                }
                
                
            }
        }
    }

    func recommendVideos(completion: @escaping ([Clip]) -> Void) {
        // Get the status URL for the recommendation
        guard let statusURL = getStatusURLFromRecommendation() else {

            print("Failed to get status URL.")
            return
        }
        // Fetch status to get the body containing video metadata
        fetchStatus(url: statusURL, textbook: "BookOfProof", subject: "math") { url in
            self.getVideos(from: url) { result in
                    print(result)
                    switch result {
                    case .success(let dataArray):
                        guard let videoUrls = dataArray as? [String] else {
                            print("Failed to parse video URLs")
                            return
                        }
                            completion(
                                videoUrls.map { url in
                                    .init(id: UUID().uuidString, videoUrl: url)
                            }
)
                    case .failure(let error):
                        print("Failed with error:", error)
                    }
                }
            }
        }
    
    func fetchStatus(url: URL, textbook: String, subject: String, completion: @escaping (String) -> Void) {
        let postData: [String: Any] = [
            "docs": [""],
            "textbook": textbook,
            "subject": subject
        ]
        
        do {
            // Serialize your data to JSON
            let jsonData = try JSONSerialization.data(withJSONObject: postData, options: [])
            // Configure the URLRequest
            var request = URLRequest(url: url)
            request.httpMethod = "POST"
            request.addValue("application/json", forHTTPHeaderField: "Content-Type")
            request.httpBody = jsonData
            
            // Perform the request
            let task = URLSession.shared.dataTask(with: request) { data, response, error in
                guard let data = data, error == nil else {
                    print("Error: \(error?.localizedDescription ?? "No error description")")
                    return
                }
                
                if let httpResponse = response as? HTTPURLResponse, httpResponse.statusCode == 200 {
                    do {
                        if let jsonObject = try JSONSerialization.jsonObject(with: data, options: []) as? [String: Any],
                           let urlValue = jsonObject["url"] as? String {
                            completion(urlValue);
                        } else {
                            print("Could not find 'url' key in the response.")
                        }
                    } catch {
                        print("Failed to parse JSON: \(error.localizedDescription)")
                    }
                    // Handle successful response
                    if let responseString = String(data: data, encoding: .utf8) {
                        print("Success: \(responseString)")
                    }
                } else {
                    // Handle server error or invalid response
                    print("Server error or invalid response")
                }
                
            }
            
            task.resume()
        } catch let error {
            print("Error serializing JSON: \(error.localizedDescription)")
        }
    }
    
    
    func getVideos(from urlString: String, completion: @escaping (Result<Any, Error>) -> Void) {
        guard let url = URL(string: urlString) else {
            completion(.failure(NSError(domain: "InvalidURL", code: -1, userInfo: nil)))
            return
        }
        
        let task = URLSession.shared.dataTask(with: url) { data, response, error in
            if let error = error {
                completion(.failure(error))
                return
            }
            
            guard let data = data else {
                completion(.failure(NSError(domain: "NoData", code: -2, userInfo: nil)))
                return
            }
            
            do {
                if let jsonObject = try JSONSerialization.jsonObject(with: data, options: []) as? [String: Any] { // Assuming the key for the array is an empty string
                    // Successfully extracted the array
                    if (jsonObject["status"] as! String != "Pending"){
                        let body = jsonObject["body"] as! [String: Any]
                        let dataArray = body[""] as! [[String: Any]]
//                        print("dataArr count is \(dataArray!.count)")
//                        print("Array count is \(self.parseS3Out(s3Out: dataArray!).count)")
                        completion(.success(self.parseS3Out(s3Out: dataArray)))
                    } else {
                        self.getVideos(from: urlString, completion: completion)
                    }
                    return
                } else {
                    completion(.failure(NSError(domain: "UnexpectedDataFormat", code: -3, userInfo: nil)))
                }
            } catch {
                completion(.failure(error))
            }
        }
        
        task.resume()
    }
    func parseS3Out(s3Out: [[String: Any]]) -> [String] {
        var out: [String] = []
        for item in s3Out { // Use a for-in loop for safer, more readable code
            if let data = item["metadata"] as? [String: Any], // Safely unwrap the metadata dictionary
               let s3VideoUri = data["s3VideoUri"] as? String { // Safely unwrap the s3VideoUri string
                out.append(s3VideoUri)
            } else {
                print("null") // This will print if either metadata is missing or s3VideoUri is not a String
            }
        }
        return out
    }
    func parseBodyForVideoUrls(_ body: [String: Any]) -> [String] {
        // Parse the body to extract video URLs
        guard let responses = body["2"] as? [[String: Any]] else {
            print("No responses found in the body.")
            return []
        }
        
        var videoUrls = [String]()
        for response in responses {
            if let metadata = response["metadata"] as? [String: Any], let s3VideoUri = metadata["s3VideoUri"] as? String {
                // Assuming s3VideoUri contains the video URL
                videoUrls.append(s3VideoUri)
            }
        }
        
        return videoUrls
    }
    
    func getStatusURLFromRecommendation() -> URL? {
        let requestBuilder = QueryCacheAPI.queryCacheGetReccomendationCreateWithRequestBuilder()
        // Ensure the URLString can be converted to a URL
        guard let url = URL(string: requestBuilder.URLString) else {
            return nil
        }

        // Optionally, create a URLRequest if you need to inspect or modify it further
        var urlRequest = URLRequest(url: url)
        urlRequest.httpMethod = requestBuilder.method

        // Apply headers from the RequestBuilder to the URLRequest
        requestBuilder.headers.forEach { header in
            urlRequest.setValue(header.value, forHTTPHeaderField: header.key)
        }

        // Now, urlRequest contains the full configuration from the RequestBuilder
        // If you just need the URL and not the fully configured URLRequest, you can return it directly:
        return url
    }
}
