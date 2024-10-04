import { getAuth } from "@clerk/nextjs/server"
import prisma from "@/lib/utils/prisma";

export async function GET(req) {
    const { userId } = getAuth(req);

    if (!userId) {
        return res.status(401).json({ error: "Unauthorized" });
    }


    const chats = await prisma.chat.findMany({
        where: {
            userId
        }
    })

    const formatted = chats.map((chat) => {
        const parsed = chat.chat.map((msg) => JSON.parse(msg));
        return {
            title: chat.title,
            chatInfo: parsed,
            id: chat.id,
            subject: chat.subject
        }
    })
    await prisma.$disconnect()
    return new Response(JSON.stringify({
        chats: formatted
    }), { status: 200 });
}

export async function POST(req, res) {
    const { userId } = getAuth(req);

    if (!userId) {
        return res.status(401).json({ error: "Unauthorized" });
    }

    const { title, chatInfo, id, subject } = await req.json();

    await prisma.chat.create({
        data: {
            userId,
            title,
            id,
            subject,
            chat: chatInfo
        }
    })

    const chats = await prisma.chat.findMany({
        where: {
            userId: userId
        }
    })

    const formatted = chats.map((chat) => {
        const parsed = chat.chat.map((msg) => JSON.parse(msg));
        return {
            title: chat.title,
            chatInfo: parsed,
            id: chat.id,
            subject: chat.subject
        }
    })
    await prisma.$disconnect()
    return new Response(JSON.stringify({
        chats: formatted
    }), { status: 200 });
}