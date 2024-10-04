//
//  SubjectCell.swift
//  EduScroll
//
//  Created by Prajwal Saokar on 1/16/24.
//

import SwiftUI


struct SubjectCell: View {
    let subject: String
    var body: some View {
        HStack(spacing: 12) {
            Image(systemName: "book")
                .resizable()
                .frame(width: 48, height: 48)
                .foregroundStyle(Color(.systemGray5))
            Text(subject)
                .font(.subheadline)
                .fontWeight(.semibold)
            Spacer()
        }
    }
}

#Preview {
    SubjectCell(subject: "Biology")
}
