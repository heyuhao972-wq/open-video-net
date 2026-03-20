const API_BASE = "http://localhost:8080"
const RECOMMEND_BASE = "http://localhost:8082"

const PLATFORM_MAP = {
    platformA: "http://localhost:8080",
    platformB: "http://localhost:8084"
}

const RECOMMEND_MAP = {
    R1: "http://localhost:8082",
    R2: "http://localhost:8086"
}

function getPlatformBase(platformId){
    return PLATFORM_MAP[platformId] || API_BASE
}

function getRecommendBase(){
    const saved = localStorage.getItem("recommend_base")
    if (saved){
        return saved
    }
    return RECOMMEND_BASE
}

function setRecommendBase(base){
    localStorage.setItem("recommend_base", base)
}

function getContentBase(){
    const saved = localStorage.getItem("content_base")
    if (saved){
        return saved
    }
    return API_BASE
}

function setContentBase(base){
    localStorage.setItem("content_base", base)
}

function parseVideoURI(uri){
    if (!uri || typeof uri !== "string"){
        return null
    }
    if (!uri.startsWith("video://")){
        return null
    }
    const rest = uri.slice("video://".length)
    const parts = rest.split("/")
    if (parts.length < 2){
        return null
    }
    return {platformId: parts[0], videoId: parts[1]}
}

async function getVideos(){

    const res = await fetch(API_BASE + "/videos")

    const data = await res.json()

    return data.videos

}

function getAuthToken(){

    return localStorage.getItem("token")

}

async function getVideoByPlatform(platformId, videoId){
    const base = getPlatformBase(platformId)
    const res = await fetch(base + "/video/" + videoId)
    return await res.json()
}

function getUserId(){

    const token = getAuthToken()
    if (token){
        const parts = token.split(".")
        if (parts.length === 3){
            try {
                const payload = JSON.parse(atob(parts[1]))
                if (payload && payload.sub){
                    return payload.sub
                }
            } catch (e) {
                // ignore
            }
        }
    }

    return localStorage.getItem("user_id")

}

async function uploadVideo(title,description,tags,file,proof){

    const form = new FormData()

    form.append("title",title)
    form.append("description",description)
    form.append("tags",tags)
    form.append("file",file)
    if (proof){
        form.append("author_signature", proof.signature)
        form.append("author_timestamp", proof.timestamp)
        form.append("video_hash", proof.videoHash)
    }

    const token = getAuthToken()

    const headers = {}
    if (token){
        headers["Authorization"] = "Bearer " + token
    }

    const res = await fetch(getContentBase() + "/upload",{
        method:"POST",
        body:form,
        headers
    })

    return await res.json()

}

async function getRecommend(user){

    const url = getRecommendBase() + "/recommend" + (user ? ("?user=" + encodeURIComponent(user)) : "")

    const res = await fetch(url)

    const data = await res.json()

    return data

}

async function reportBehavior(videoId,type){

    const userId = getUserId()
    const token = getAuthToken()
    if (!userId || !token){
        return {error:"login required"}
    }

    const res = await fetch(getRecommendBase() + "/behavior",{
        method:"POST",
        headers:{
            "Content-Type":"application/json",
            "Authorization":"Bearer " + token
        },
        body: JSON.stringify({
            user_id: userId,
            video_id: videoId,
            type
        })
    })

    return await res.json()

}

async function searchVideos(q){

    const base = getContentBase()
    const res = await fetch(base + "/videos")
    const data = await res.json()
    const list = data.videos || []
    const query = (q || "").toLowerCase()
    if (!query){
        return list
    }
    return list.filter(v=>{
        const title = (v.title || "").toLowerCase()
        const desc = (v.description || "").toLowerCase()
        const tags = Array.isArray(v.tags) ? v.tags.join(",").toLowerCase() : ""
        return title.includes(query) || desc.includes(query) || tags.includes(query)
    })

}

async function registerUser(username,password){

    const res = await fetch(API_BASE + "/register",{
        method:"POST",
        headers:{"Content-Type":"application/json"},
        body: JSON.stringify({public_key: username})
    })

    return await res.json()

}

function b64ToBytes(b64){
    const bin = atob(b64)
    const bytes = new Uint8Array(bin.length)
    for (let i=0;i<bin.length;i++){
        bytes[i]=bin.charCodeAt(i)
    }
    return bytes
}

function bytesToB64(bytes){
    let bin = ""
    for (let i=0;i<bytes.length;i++){
        bin += String.fromCharCode(bytes[i])
    }
    return btoa(bin)
}

async function generateKeypair(){
    const keypair = await crypto.subtle.generateKey(
        {name:"Ed25519"},
        true,
        ["sign","verify"]
    )
    const publicRaw = new Uint8Array(await crypto.subtle.exportKey("raw", keypair.publicKey))
    const privatePkcs8 = new Uint8Array(await crypto.subtle.exportKey("pkcs8", keypair.privateKey))
    return {
        publicKey: bytesToB64(publicRaw),
        privateKey: bytesToB64(privatePkcs8)
    }
}

async function hashFile(file){
    const buf = await file.arrayBuffer()
    const digest = await crypto.subtle.digest("SHA-256", buf)
    const hashBytes = new Uint8Array(digest)
    let hex = ""
    for (const b of hashBytes){
        hex += b.toString(16).padStart(2,"0")
    }
    return hex
}

function buildProofMessage(videoHash, timestamp, publicKey){
    return new TextEncoder().encode(videoHash + "|" + timestamp + "|" + publicKey)
}

async function signNonce(nonceB64, privateKeyB64){
    const keyBytes = b64ToBytes(privateKeyB64)
    const key = await crypto.subtle.importKey(
        "pkcs8",
        keyBytes,
        {name:"Ed25519"},
        false,
        ["sign"]
    )
    const nonceBytes = b64ToBytes(nonceB64)
    const sig = await crypto.subtle.sign("Ed25519", key, nonceBytes)
    return bytesToB64(new Uint8Array(sig))
}

async function signMessage(messageBytes, privateKeyB64){
    const keyBytes = b64ToBytes(privateKeyB64)
    const key = await crypto.subtle.importKey(
        "pkcs8",
        keyBytes,
        {name:"Ed25519"},
        false,
        ["sign"]
    )
    const sig = await crypto.subtle.sign("Ed25519", key, messageBytes)
    return bytesToB64(new Uint8Array(sig))
}

async function loginUser(username){

    const challengeRes = await fetch(API_BASE + "/login/challenge",{
        method:"POST",
        headers:{"Content-Type":"application/json"},
        body: JSON.stringify({public_key: username})
    })

    const challenge = await challengeRes.json()
    if (challenge.error){
        return challenge
    }
    const nonce = challenge.nonce
    const privateKey = localStorage.getItem("private_key")
    if (!privateKey){
        return {error:"private key missing, please register again"}
    }

    const signature = await signNonce(nonce, privateKey)

    const res = await fetch(API_BASE + "/login",{
        method:"POST",
        headers:{"Content-Type":"application/json"},
        body: JSON.stringify({public_key: username,nonce,signature})
    })

    return await res.json()

}

async function followAuthor(authorId){

    const token = getAuthToken()
    if (!token){
        return {error:"login required"}
    }

    const res = await fetch(getRecommendBase() + "/follow",{
        method:"POST",
        headers:{
            "Content-Type":"application/json",
            "Authorization":"Bearer " + token
        },
        body: JSON.stringify({author_id: authorId})
    })

    return await res.json()

}

async function unfollowAuthor(authorId){

    const token = getAuthToken()
    if (!token){
        return {error:"login required"}
    }

    const res = await fetch(getRecommendBase() + "/unfollow",{
        method:"POST",
        headers:{
            "Content-Type":"application/json",
            "Authorization":"Bearer " + token
        },
        body: JSON.stringify({author_id: authorId})
    })

    return await res.json()

}
