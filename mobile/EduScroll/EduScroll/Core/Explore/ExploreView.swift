import SwiftUI

struct ExploreView: View {
    @StateObject var viewModel = SubjectViewModel()

    var body: some View {
        NavigationStack {
            if viewModel.isLoading {
                ProgressView("Loading")
            } else {
                ScrollView {
                    LazyVStack {
                        ForEach(viewModel.subjects, id: \.self) { subject in
                            SubjectCell(subject: subject)
                                .padding(.horizontal)
                        }
                    }
                }
                .navigationTitle("Subjects")
                .navigationBarTitleDisplayMode(.inline)
                .padding(.top)
            }
        }
    }
}

#if DEBUG
struct ExploreView_Previews: PreviewProvider {
    static var previews: some View {
        ExploreView()
    }
}
#endif
