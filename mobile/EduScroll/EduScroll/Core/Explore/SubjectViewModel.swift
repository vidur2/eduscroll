import Foundation

class SubjectViewModel: ObservableObject {
    @Published var subjects: [String] = []
    @Published var isLoading = false
    
    init() {
        fetchSubjects()
    }

    func fetchSubjects() {
        isLoading = true
    
        guard let basePath = Optional(SwaggerClientAPI.basePath),
        let url = URL(string: basePath + "/subjects")
                
        else {
               print("Invalid URL")
               return
           }

        var request = URLRequest(url: url)
        request.httpMethod = "GET"

        URLSession.shared.dataTask(with: request) { data, response, error in
            DispatchQueue.main.async {
                self.isLoading = false
            }
            
            guard let data = data else {
                print("No data received: \(error?.localizedDescription ?? "Unknown error")")
                return
            }

            do {
                let jsonObject = try JSONSerialization.jsonObject(with: data, options: [])
                
                if let jsonDictionary = jsonObject as? [String: Any], let collections = jsonDictionary["collections"] as? [String] {
                    DispatchQueue.main.async {
                        self.subjects = collections
                    }
                } else {
                    print("Invalid JSON format")
                }
                print("Subjects are ", self.subjects)
            } catch {
                print("JSON parsing error: \(error.localizedDescription)")
            }
        }.resume()
    }

}
