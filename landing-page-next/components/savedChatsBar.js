'use client'

import { useEffect, useState } from "react"
import styles from "../app/page.module.css"
import { Modal, ModalClose, Typography, ModalDialog } from "@mui/joy"
import { LinearGradient } from "react-text-gradients"
import Select from "@mui/joy/Select"
import Option from "@mui/joy/Option"
import { NewChatDisplay } from "./newChatDisplay"

/**
 * 
 * @param {
 *   []savedChats {
 *     title: str,
 *     id: str,
 *     chatInfo: []{question: str, response: str}
 * }
 * } param0 
 * @returns 
 */
export function SavedChatsBar({ savedChats, setSavedChats, subjects, chatId, screenWidth }) {
    const [open, setOpen] = useState(false)
    const [newChatTitle, setNewChatTitle] = useState("");
    const [subject, setSubject] = useState("");
    const handleModalOpen = () => {
        setOpen(true);
    }

    const handleSubjectChange = (e, nv) => {
        setSubject(nv);
    }

    const handleModalClose = () => {
        setOpen(false);
    }

    const handleTextChange = (e) => {
        setNewChatTitle(e.target.value);
    }

    const handleModalSubmit = () => {
        const newChat = {
            title: newChatTitle,
            chatInfo: new Array(),
            id: crypto.randomUUID(),
            subject
        }
        console.log(savedChats)
        fetch(("/api/db/chat"), {
            method: "POST",
            body: JSON.stringify(newChat)
        }).then((res) => res.json().then((body) => {
            if (res.status == 200) {
                setSavedChats(body.chats)
            }
            setOpen(false);
        }));
    }

    function getTextWidth(text, font) {
        const canvas = getTextWidth.canvas || (getTextWidth.canvas = document.createElement("canvas"));
        const context = canvas.getContext("2d");
        context.font = font;
        const metrics = context.measureText(text);
        return metrics.width;
    }

    return (
        <div style={{display: "flex", flexDirection: "column", maxHeight: "inherit", whiteSpace: "nowrap"}}>
            <div>
                <h2 style={{paddingLeft: "8%" }}><a href={"/chat"} style={{textDecoration: "none", color: "black"}}><LinearGradient gradient={['to left', '#3395FF ,#3358ff']}>Previous</LinearGradient> Chats</a></h2>
                {savedChats.toReversed().map((chat, i) => 
                    <div key={`chatBarDiv_${i}`} style={{marginTop: "8%", marginBottom: "8%", textOverflow: "hidden", maxWidth: "200px" }}>
                        {
                            chat.id == chatId ? (
                                <a href={`/chat/${chat.id}`} key={`chatBar_${i}`} style={{textDecoration: "none", paddingRight: `${200 - getTextWidth(chat.title)}px`, paddingLeft: "8%", paddingBottom: "4%", paddingTop: "4%", borderRadius: "5px", backgroundColor: "#15406e", color: "white"}} className={styles.chatBarStyle}>{chat.title}</a>
                            ) : (
                                <a href={`/chat/${chat.id}`} key={`chatBar_${i}`} style={{textDecoration: "none", paddingRight: `${200 - getTextWidth(chat.title)}px`, paddingLeft: "8%", paddingBottom: "4%", paddingTop: "4%", borderRadius: "5px"}} className={styles.chatBarStyle}>{chat.title}</a>
                            )
                        }
                    </div>
                )}
            </div>
            <div style={{position: "absolute", bottom: 5}}>
                <button style={{ float: "bottom", marginLeft: "2%", whiteSpace: "nowrap", paddingTop: "5px", paddingBottom: "5px"}} onClick={handleModalOpen}>New Chat</button>
            </div>
            <NewChatDisplay subjects={subjects} open={open} handleSubjectChange={handleSubjectChange} handleModalClose={handleModalClose} handleTextChange={handleTextChange} handleModalSubmit={handleModalSubmit} />
        </div>
    )
}