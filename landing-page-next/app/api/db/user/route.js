import { getAuth } from "@clerk/nextjs/server";
import prisma from "@/lib/utils/prisma";


export async function GET(req) {
    const { userId } = getAuth(req);
    if (!userId) {
        return res.status(401).json({ error: "Unauthorized" });
    }

    const { subjects, textbooks } = await prisma.user.findUnique({
        where: {
            clerkId: userId
        }
    })

    await prisma.$disconnect();

    return new Response(JSON.stringify({
        userInfo: {
            subjects,
            textbooks
        }
    }), { status: 200 })
}


export async function POST(req) {
    const { userId } = getAuth(req);
    if (!userId) {
        return res.status(401).json({ error: "Unauthorized" });
    }

    const { subjects, textbooks } = await req.json();

    const user = await prisma.user.update({
        data: {
            subjects: subjects,
            textbooks: textbooks
        }
    });

    await prisma.$disconnect();
    
    return new Response(JSON.stringify({
        userInfo: {
            subjects: user.subjects,
            textbooks: user.textbooks
        }
    }), { status: 200 });
}