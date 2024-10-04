import { Webhook } from 'svix'
import { headers } from 'next/headers'
import prisma from '@/lib/utils/prisma'
 
export async function POST(req) {
  // You can find this in the Clerk Dashboard -> Webhooks -> choose the webhook
    const WEBHOOK_SECRET = process.env.WEBHOOK_SECRET
    
    if (!WEBHOOK_SECRET) {
        throw new Error('Please add WEBHOOK_SECRET from Clerk Dashboard to .env or .env.local')
    }
    
  // Get the headers
    const headerPayload = headers();
    const svix_id = headerPayload.get("svix-id");
    const svix_timestamp = headerPayload.get("svix-timestamp");
    const svix_signature = headerPayload.get("svix-signature");
  // If there are no headers, error out
    if (!svix_id || !svix_timestamp || !svix_signature) {
        return new Response('Error occured -- no svix headers', {
                status: 400
            })
    }

    // Get the body
    const payload = await req.json()
    const body = JSON.stringify(payload);

    // Create a new Svix instance with your secret.
    const wh = new Webhook(WEBHOOK_SECRET);

    let evt;

    // Verify the payload with the headers
    try {
        evt = wh.verify(body, {
            "svix-id": svix_id,
            "svix-timestamp": svix_timestamp,
            "svix-signature": svix_signature,
            })
    } catch (err) {
    console.error('Error verifying webhook:', err);
        return new Response('Error occured', {
                status: 400
            })      
    }

  // Get the ID and type
    const { id } = evt.data;
    const eventType = evt.type;
    console.log(eventType);
    if (evt.type === 'user.created' || evt.type === 'user.updated') {
        const firstName = evt.data.first_name;
        const emails = evt.data.email_addresses.map((email) => email.email_address);  
        console.log(emails)
        const lastName = evt.data.last_name;
        const clerkId = evt.data.id;

        await prisma.user.upsert({
          where: {
            clerkId
          },
          update: {
            firstName,
            lastName,
            emails,
          },
          create: {
            firstName,
            lastName,
            emails,
            clerkId
          }
        })
        await prisma.$disconnect()
    } else if (evt.type === 'user.deleted') {
      await prisma.user.delete({
        where: {
          clerkId: evt.data.id
        }
      })
      await prisma.$disconnect()
    }

    return new Response('', { status: 200 })
}