//
//  MainTabView.swift
//  EduScroll
//
//  Created by EduScroll Development
//

import SwiftUI

struct MainTabView: View {
    var body: some View {
        TabView {
            Text("Profile")
                .tabItem {
                    VStack {
                        Image(systemName: "person")
                        Text("Profile")
                    
                    }
                }
            FeedView()
                .tabItem {
                    VStack {
                        Image(systemName: "house")
                        Text("Home")
                    }
                }
            ExploreView()
                .tabItem {
                    VStack {
                        Image(systemName: "book")
                        Text("Explore")
                    }
                }
        }
    }
}
                             
#Preview {
    MainTabView()
    }
