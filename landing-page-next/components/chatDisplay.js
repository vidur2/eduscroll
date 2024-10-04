import { CircularProgress } from "@mui/joy"

export function ChatDisplay({ chat, loading, query }) {
    return (
        <div style={{overflow: "scroll", marginLeft: "10%", height: "100%", display: "flex", flexDirection: "column-reverse", maxHeight: "inherit" }}>
            {
                loading ? (
                    <div key={`chatDisp_${chat.length}`} >
                        <h3>You:</h3>
                        <p style={{fontSize: "15px"}}>{query}</p>
                        <h3>Response: </h3>
                        <div style={{width: "50%", marginLeft: "30%"}}>
                        <CircularProgress size="lg" />
                        </div>
                    </div>
                ) : (
                    <div></div>
                )
            }
            {
                chat.toReversed().map((chat, i) => (
                    <div key={`chatDisp_${i}`} >
                        <h3>You:</h3>
                        <p style={{fontSize: "15px"}}>{chat.question}</p>
                        <h3>Response: </h3>
                        <p style={{fontSize: "15px"}}>{chat.response}</p>
                    </div>
                ))
            }
        </div>
    )
}