'use client'

import Image from 'next/image';
import submitIco from "../public/up-arrow.svg";
import styles from "../app/page.module.css";
import Select from '@mui/joy/Select';
import Option from '@mui/joy/Option';

export function MessageBar({ msgSubmitHandler, chatHandler }) {
    
    return (
        <div style={{borderWidth: "1px", borderColor: "black", borderStyle: "solid", borderRadius: "5px", marginRight: "25%", bottom: "1%", position: "fixed", left: "25%", right: 0 }}>
        <form id="messageForm" onSubmit={msgSubmitHandler} >

            <input type="text" style={{borderColor: "transparent", borderWidth: "0", width: "93%", fontSize: "100%", marginLeft: "2%", outline: "none", marginBottom: "1%" }} id="msgInput" onChange={chatHandler} ></input>
            <button type="submit" style={{paddingLeft: "0", paddingRight: "0%", borderWidth: "0px", borderColor: "transparent", marginTop: "1%" }}
            ><
                Image
                src={submitIco} 
                width={15}
            />

            </button>
        </form>
        </div>
    )
}