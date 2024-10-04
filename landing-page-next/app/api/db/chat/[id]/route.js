import { getAuth } from "@clerk/nextjs/server"
import { PrismaClient } from "@prisma/client";

export async function POST(req) {
    const splitUrl = req.url.split("/");
    const id = splitUrl[splitUrl.length - 1]
    const { userId } = getAuth(req);

    if (!userId) {
        return new Response({}, {status: 400})
    }

    const { chatInfo } = await req.json();
    const prisma = new PrismaClient();

    const tmp = await prisma.chat.findUnique({
        where: {
            id
        }
    })

    const response = "change";    // Call langchain here maybe
    chatInfo["response"] = response;

    tmp.chat.push(JSON.stringify(chatInfo));

    const chat = await prisma.chat.update({
        where: {
            id
        },
        data: {
            chat: tmp.chat
        }
    })

    return new Response(JSON.stringify({
        chatInfo: chat.chat.map((msg) => JSON.parse(msg))
    }), { status: 200 })
}

const queryBackend = (prompt, subject) => {

}