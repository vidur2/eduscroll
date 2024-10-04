import { LinearGradient } from "react-text-gradients"
export function Bar() {
    return (
        <div style={{ backgroundColor: "white", paddingTop: "2%", paddingBottom: "1.5%"}}>
            <div style={{ marginLeft: "3%" }}>
                <a href="#download" className="linkStyle">Download</a>
                <a href="#technology" className="linkStyle">The Technology</a>
                <a href="#contactSales" className="linkStyle">Contact Sales</a>
            </div>
            <div style={{paddingTop: "1%" }}></div>
            <div style={{ textAlign: "center"}}>
            <h1 style={{ color: "#3395FF", display:"inline", fontSize: "500%" }}><LinearGradient gradient={['to left', '#3395FF ,#3358ff']}>Edu</LinearGradient></h1><h1 style={{ color: "#000000",display:"inline", fontSize: "500%" }}>Scroll</h1>
            </div>
        </div>
    )
}