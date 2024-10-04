'use client'

import { useState } from "react"
import styles from "../app/page.module.css"
import { Modal, ModalClose, Typography, ModalDialog, Sheet } from "@mui/joy"
import { LinearGradient } from "react-text-gradients"
import { useRouter } from "next/navigation"
import { UserButton } from "@clerk/nextjs"
import Select from "@mui/joy/Select"
import Option from "@mui/joy/Option"
import {NewChatDisplay} from "./newChatDisplay"

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
export function SavedChatsBarFull({ savedChats, setSavedChats, subjects }) {
    const [open, setOpen] = useState(false)
    const router = useRouter();
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
        const uuid = crypto.randomUUID();
        const newChat = {
            title: newChatTitle,
            chatInfo: new Array(),
            id: uuid,
            subject
        }
        console.log(savedChats)
        fetch(("/api/db/chat"), {
            method: "POST",
            body: JSON.stringify(newChat)
        }).then((res) => {
            if (res.status == 200) {
                router.push(`/chat/${uuid}`)
            } else {
                setOpen(false);
            }
        });
    }

    return (
        <div style={{maxHeight: "50px" }}>
            <div style={{textAlign: "center"}}>

                <h1><LinearGradient gradient={['to left', '#3395FF ,#3358ff']} style={{marginBottom: "10%"}}>Previous</LinearGradient> Chats</h1>
                {savedChats.toReversed().map((chat, i) => 
                    <div key={`chatBarDiv_${i}`} style={{marginTop: "2%", marginBottom: "2%"}}>
                        <a href={`/chat/${chat.id}`} key={`chatBar_${i}`} style={{textDecoration: "none", paddingRight: "32%", paddingLeft: "30%", paddingBottom: "1%", paddingTop: "1%", borderRadius: "5px"}} className={styles.chatBarStyle}>{chat.title}</a>
                    </div>
                )}
                <button onClick={handleModalOpen} style={{marginTop: "2%", paddingTop: "5px", paddingBottom: "5px"}}>New Chat</button>
            </div>
            <NewChatDisplay subjects={subjects} open={open} handleSubjectChange={handleSubjectChange} handleModalClose={handleModalClose} handleTextChange={handleTextChange} handleModalSubmit={handleModalSubmit} />
            {/* <Modal open={open}>
            <ModalDialog style={{ textAlign: "center" }}>
                <Typography>Enter Chat Title</Typography>
                    <input type="text" onChange={handleTextChange} style={{ borderWidth: "1px", borderColor: "black", borderRadius: "5px", height: "150%" }}></input>
                    <Select onChange={handleSubjectChange}>
                        {subjects.map((subject, i) => (
                            <Option value={subject} key={`dropdown${i}`}>
                                {subject}
                            </Option>
                        ))}
                    </Select>
                    <table>
                    <tbody>
                    <td>
                        <button onClick={handleModalClose} className={styles.cancel}>Cancel</button>
                    </td>
                    <td>
                        <button onClick={handleModalSubmit}>Done</button>
                    </td>
                    </tbody>
                    </table>
            </ModalDialog>
            </Modal> */}
        </div>
    )
}