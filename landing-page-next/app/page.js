'use client'

import { Bar } from '@/components/bar';
import {LinearGradient} from "react-text-gradients"; 
import Image from 'next/image';
import { useRouter } from 'next/navigation'
import Head from 'next/head';

const calendlyEmbed = "<!-- Calendly badge widget begin --><link href='https://assets.calendly.com/assets/external/widget.css' rel='stylesheet'><script src='https://assets.calendly.com/assets/external/widget.js' type='text/javascript' async></script><script type='text/javascript'>window.onload = function() { Calendly.initBadgeWidget({ url: 'https://calendly.com/vmod2005', text: 'Schedule time with me', color: '#0069ff', textColor: '#ffffff', branding: true }); }</script><!-- Calendly badge widget end -->";

function App() {

  return (
    <div style={{ backgroundColor: "#f1f1f1" }}>
      <Head>
        <title>Eduscroll</title>
      </Head>
      <Bar></Bar>
      <div style={{ padding: "2%", alignItems: "center"}}>
        <img src="/edutokBody.png" width="85%" style={{display: "block", marginLeft: "auto", marginRight: "auto", borderRadius: "24px"}}></img>
      </div>
      <div style={{ backgroundColor: "white" }}>
      <h1 style={{ textAlign: "center", paddingTop: "5%", fontSize: "300%", fontWeight: "550", marginBottom: "8%" }} id="features">
        Features
      </h1>
      <table id="download" style={{ marginLeft: "20%" }}>
      <td style={{ paddingRight: "5%" }}>
        <img src="/edutokPhoneImage.png" width= "90%" style={{display: "block", marginLeft: "auto", marginRight: "auto", borderRadius: "24px"}}></img>

        </td>
        <td style={{ position: "relative", bottom: "300px", paddingLeft: "5%" }}>
          <h2 style={{ fontWeight: "550", fontSize: "150%" }}>
            Learning.  <LinearGradient gradient={['to left', '#3395FF ,#3358ff']}>Reimagined.</LinearGradient>
          </h2>
          <p style={{ maxWidth: "50%" }}>Bite-sized content pulled directly out of your selected textbooks, catered for your K-12 education needs. This is education for the modern era. At your fingertips.</p>
          <button style={{paddingTop: "5px", paddingBottom: "5px", paddingLeft: "15%", paddingRight: "15%"}}>Try it out</button>
        </td>
      </table>
      <table style={{marginLeft: "20%", marginRight: "20%"}} id="technology">
        <td style={{ position: "relative", bottom: "350px" }}>
        <h2 style={{ fontWeight: "550", fontSize: "150%", whiteSpace: "nowrap" }}>
          <p style={{ marginBottom: "0px" }}> See how Generative <LinearGradient gradient={['to left', '#3395FF ,#3358ff']}>AI</LinearGradient></p> can power learning
          </h2>
          <p style={{ maxWidth: "50%" }}>EduTok utilizes generative AI to create a personalized learning experience for your student. We utilize the power of LLMs, image generation, and text to speech to generate engaging videos. Every time.</p>
          <button style={{paddingTop: "5px", paddingBottom: "5px", paddingLeft: "15%", paddingRight: "15%"}}>Learn More</button>
        </td>
        <td>
          <img src="/edutokAIColorful.png" width= "150%"  style={{display: "block", marginLeft: "auto", marginRight: "auto", borderRadius: "24px"}}></img>
        </td>
      </table>
      <table style={{ marginLeft: "15%", marginTop: "-5%" }}>
        <td style={{ paddingRight: "5%" }}>
          <img src="/districtLevel.png" width= "77.5%" layout='fill' style={{display: "block", marginLeft: "auto", marginRight: "auto", borderRadius: "24px"}}></img>
        </td>
        <td style={{ textAlign: "justify", position: "relative", bottom: "300px", marginLeft: "5%" }} id="contactSales">
          <h2 style={{ fontWeight: "550", fontSize: "150%" }}>
            Service at <LinearGradient gradient={['to left', '#3395FF ,#3358ff']}>every</LinearGradient> level
          </h2>
          <p style={{ maxWidth: "60%" }}>For our larger-scale enterprise products, get in touch with our sales team so we can tailor our service to your needs</p>
          <button style={{paddingTop: "5px", paddingBottom: "5px", paddingLeft: "15%", paddingRight: "15%"}}>Contact us</button>
        </td>
      </table>
      </div>
    </div>
  );
}

export default App;
