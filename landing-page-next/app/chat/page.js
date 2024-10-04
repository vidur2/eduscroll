'use client'
import { SavedChatsBarFull } from "@/components/savedChatsFull";
import { useState, useEffect } from "react";
import { UserButton } from "@clerk/nextjs";
import { Stack } from "@mui/joy";
import { styled } from '@mui/joy/styles';

function Chat() {
    const [savedChats, setSavedChats] = useState([])
    const [subjects, setSubjects] = useState([]);

    useEffect(() => {
        fetch("/api/db/chat").then((res) => {
            if (res.status == 200) {
                res.json().then((body) => setSavedChats(body.chats))
            } else {
                console.log("chat fail")
            }
        })

        fetch("/api/db/user").then((res) => {
            if (res.status == 200) {
                res.json().then((body) => setSubjects(body.userInfo.subjects))
            } else {
                console.log("user fail")
            }
        })
    }, [])

    return (
        <div>
            <div style={{float: "right", marginRight: "5%"}}>
            <UserButton/>
            </div>
        <div style={{marginRight: "5%"}}>
        <div>
        <SavedChatsBarFull savedChats={savedChats} subjects={subjects} />
        </div>
        </div>
        </div>
    )
}

export default Chat;