'use client'

import { MessageBar } from "@/components/messageBar";
import { Grid } from "@mui/joy";
import { SavedChatsBar } from "@/components/savedChatsBar";
import { ChatDisplay } from "@/components/chatDisplay";
import { useEffect, useState } from "react";
import { UserButton } from "@clerk/nextjs";
import Select from '@mui/joy/Select';
import Option from '@mui/joy/Option';

function Chat() {
    const [chat, setChat] = useState([])
    const [savedChats, setSavedChats] = useState([])
    const [msg, setMsg] = useState("")
    const [vidLoading, setVidLoading] = useState(false);
    const [textbooks, setTextbooks] = useState([]);
    const [subjects, setSubjects] = useState([]);
    const [uuid, setUuid] = useState("");

    const msgSubmitHandler = (e) => {
        e.preventDefault();
        const splitPath = window.location.pathname.split("/");
        const slug = splitPath[splitPath.length - 1];
        document.getElementById("msgInput").value = "";
        setVidLoading(true);
        fetch(`/api/db/chat/${slug}`, {
            method: "POST",
            body: JSON.stringify({
                chatInfo: { question: msg }
            })
        }).then((res) => {
            if (res.status == 200) {
                res.json().then((body) => {
                    setVidLoading(false);
                    setChat(body.chatInfo);
                })
            }
        })
    }

    const chatHandler = (e) => {
        e.preventDefault();
        setMsg(e.target.value);
    }

    useEffect(() => {
        setUuid(window.location.href.split("/")[4]);
        const splitPath = window.location.pathname.split("/");
        const slug = splitPath[splitPath.length - 1];
        
        fetch("/api/db/chat").then(async(res) => {
            if (res.status == 200) {
                const body = await res.json();
                setSavedChats(body.chats)
                const chatTmp = body.chats.filter((chat) => chat.id === slug);
                if (chatTmp.length > 0) {
                    setChat(chatTmp[0].chatInfo);
                    console.log(chatTmp[0].chatInfo)
                } else {
                    console.log("failed to load chat properly")
                }
            } else {
                console.log("chat fail")
            }
        })

        fetch("/api/db/user").then((res) => {
            if (res.status == 200) {
                res.json().then((body) => {
                    setSubjects(body.userInfo.subjects);
                    setTextbooks(body.userInfo.textbooks)
                })
            } else {
                console.log("user fail")
            }
        })
    }, [])


    return (
        <div>
            <div style={{float: "right", marginRight: "5%", marginTop: "1%" }}>
                <table>
                    <Select placeholder="Choose textbook" >
                    {
                        textbooks.map((textbook, i) => {
                            <Option key={`textbookMap${i}`} value={textbook}>{textbook}</Option>
                        })
                    }
                    </Select>
                </table>
            </div>
            <Grid container spacing={0} sx={{ flexGrow: 0 }}>
            <Grid xs={2}>
            <div style={{display: "inline-block", maxHeight: "50px", width: "100%"}}>
            <SavedChatsBar savedChats={savedChats} setSavedChats={setSavedChats} subjects={subjects} chatId={uuid} screenWidth={window.screen.width} />
            </div>
            </Grid>
            <Grid xs={10}>
            <div style={{display: "inline-block", maxHeight: "700px", width: "100%"}}>
            <ChatDisplay chat={chat} loading={vidLoading} query={msg} />
            </div>
            <MessageBar style={{verticalAlign: "bottom"}} msgSubmitHandler={msgSubmitHandler} chatHandler={chatHandler} ></MessageBar>
            </Grid>
            </Grid>
        </div>
    )
}

export default Chat;