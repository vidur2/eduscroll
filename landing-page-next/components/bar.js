'use client'

import { LinearGradient } from "react-text-gradients"
import styles from "../app/page.module.css"
import { SignInButton, SignedIn, SignedOut, UserButton } from "@clerk/nextjs"

export function Bar() {
    return (
        <div style={{ backgroundColor: "white", paddingTop: "2%", paddingBottom: "1.5%" }}>
            <div style={{ marginLeft: "3%" }}>
                <SignedIn>
                    <a href="/chat" className={styles.linkStyle}>Chat</a>
                </SignedIn>
                <a href="#download" className={styles.linkStyle}>Download</a>
                <a href="#technology" className={styles.linkStyle}>The Technology</a>
                <a href="#contactSales" className={styles.linkStyle}>Contact Sales</a>
                <span style={{float: "right", marginRight: "5%" }}>
                    <SignedOut>
                        <SignInButton style={{ paddingLeft: "10%", paddingRight: "10%", backgroundColor: "", whiteSpace: "nowrap", paddingTop: "5px", paddingBottom: "5px" }}/>
                    </SignedOut>
                    <SignedIn>
                        <UserButton className={styles.userButton}></UserButton>
                    </SignedIn>
                </span>
            </div>
            <div style={{paddingTop: "1%" }}></div>
            <div style={{ textAlign: "center"}}>
            <h1 style={{ color: "#3395FF", display:"inline", fontSize: "500%" }}><LinearGradient gradient={['to left', '#3395FF ,#3358ff']}>Edu</LinearGradient></h1><h1 style={{ color: "#000000",display:"inline", fontSize: "500%" }}>Scroll</h1>
            </div>
        </div>
    )
}