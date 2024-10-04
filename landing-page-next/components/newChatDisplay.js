import { Modal, Typography, ModalDialog, Select, Option } from "@mui/joy"
import styles from "../app/page.module.css"
import { useState } from "react";

export function NewChatDisplay({ subjects, open, handleSubjectChange, handleModalClose, handleTextChange, handleModalSubmit }) {
    return (
        <Modal open={open}>
            <ModalDialog style={{ textAlign: "center" }}>
                <Typography>Enter Chat Title</Typography>
                    <input type="text" onChange={handleTextChange} style={{ borderWidth: "1px", borderColor: "black", borderRadius: "5px", height: "150%" }}></input>
                    <Select onChange={handleSubjectChange} key={`dropdown`}>
                        {subjects.map((subject, i) => (
                            <Option value={subject} key={`optionMapSavedChatbar_${i}`}>
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
            </Modal>
    )
}