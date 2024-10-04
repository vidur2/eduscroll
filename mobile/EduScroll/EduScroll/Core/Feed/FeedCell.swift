//
//  FeedCell.swift
//  EduScroll
//
//  Created by EduScroll Development 
//

import SwiftUI
import AVKit

struct FeedCell: View {
    let clip: Clip
    var body: some View {
        ZStack {
            VideoPlayer(player: AVPlayer(url: URL(string: clip.videoUrl)!))
                .containerRelativeFrame([.horizontal, .vertical])
            VStack {
                Spacer()
            }
        }
    }
}


#Preview {
    FeedCell(clip: Clip(id: "CEEFD078-8080-422C-BAE6-88227810D9E5", videoUrl: "https://eduscroll-video-output.s3.amazonaws.com/3d5e2eba-af80-4b7e-a573-145b8fc53495-combined.mp4?AWSAccessKeyId=AKIAQHVKN6JGMZ2VCT7F&Signature=t1GJ2lTDsEaSTN6Y1%2BdLe4HrE4c%3D&Expires=1708722218"))
}
