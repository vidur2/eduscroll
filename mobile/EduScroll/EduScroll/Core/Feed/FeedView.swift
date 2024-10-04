import SwiftUI

struct FeedView: View {
    @StateObject var viewModel = FeedViewModel()

    var body: some View {
        if viewModel.isLoading {
            ProgressView("Loading")
        } else {
            ScrollView {
                LazyVStack(spacing: 0) {
                    ForEach(viewModel.clips) {clip in
                        FeedCell(clip: clip)                        
                    }
                } .scrollTargetLayout()
            }.scrollTargetBehavior(.paging)
                .ignoresSafeArea()
            
        }
    }
}

#Preview {
    FeedView()
}
