import { Modal, ModalDialog, Typography } from "@mui/joy";
import { useRouter } from "next/navigation";
import { useState } from "react";
import styles from "../app/page.module.css";

export function ChatDisplayReplacement({ savedChats, setSavedChats }) {
    const router = useRouter();
    const [open, setOpen] = useState(false);
    const [newChatTitle, setNewChatTitle] = useState("");
    const handleModalOpen = () => {
        setOpen(true);
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
            id: uuid
        };
        console.log(savedChats)
        setSavedChats([...savedChats, newChat]) // POST to backend
        setOpen(false);
        router.push(`/chat/${uuid}`)
    }
    return (
        <div>
        <div style={{position: "absolute", bottom: 5}}>
                <button style={{ float: "bottom", marginLeft: "2%", whiteSpace: "nowrap" }} onClick={handleModalOpen}>New Chat</button>
            </div>
            <Modal open={open}>
            <ModalDialog style={{ textAlign: "center" }}>
                <Typography>Enter Chat Title</Typography>
                    <input type="text" onChange={handleTextChange} style={{ borderWidth: "1px", borderColor: "black", borderRadius: "5px", height: "150%" }}></input>
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
            </Modal>
        </div>
    )
}