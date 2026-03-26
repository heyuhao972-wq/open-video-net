const API_BASE = "http://localhost:9000"
const STREAM_BASE = "http://localhost:8081"
const STORAGE_BASE = "http://localhost:9001"

const PLATFORM_MAP = {
    platformA: "http://localhost:9000",
    platformB: "http://localhost:8084"
}

const RECOMMEND_OPTIONS = [
    {id:"recA", label:"鎺ㄨ崘A", url:"http://localhost:9002"}
]

function getPlatformBase(platformId){
    return PLATFORM_MAP[platformId] || API_BASE
}

function getRecommendOptions(){
    return RECOMMEND_OPTIONS.slice()
}

function normalizeRecommendBase(url){
    if (!url){
        return ""
    }
    if (url.indexOf("localhost:8082") >= 0){
        return "http://localhost:9002"
    }
    return url
}

function getRecommendBase(){
    const stored = localStorage.getItem("recommend_base")
    const normalized = normalizeRecommendBase(stored)
    if (normalized && normalized !== stored){
        localStorage.setItem("recommend_base", normalized)
    }
    return normalized || RECOMMEND_OPTIONS[0].url
}

function setRecommendBase(url){
    localStorage.setItem("recommend_base", url)
}

function getStorageBase(){
    return localStorage.getItem("storage_base") || STORAGE_BASE
}

async function parseJsonSafe(res){
    const text = await res.text()
    if (!text){
        return {}
    }
    try {
        return JSON.parse(text)
    } catch (e) {
        return {error:"invalid json", raw:text}
    }
}

async function getMyProfile(){
    const token = getAuthToken()
    if (!token){
        return {error:"login required"}
    }
    const res = await fetch(API_BASE + "/me/profile",{
        headers:{
            "Authorization":"Bearer " + token
        }
    })
    if (res.ok){
        return await parseJsonSafe(res)
    }
    const userId = getUserId()
    if (userId){
        const fallback = await fetch(API_BASE + "/profile/" + encodeURIComponent(userId))
        if (fallback.ok){
            return await parseJsonSafe(fallback)
        }
    }
    return await parseJsonSafe(res)
}

async function updateProfile(nickname, avatarUrl, bio){
    const token = getAuthToken()
    if (!token){
        return {error:"login required"}
    }
    const res = await fetch(API_BASE + "/profile",{
        method:"POST",
        headers:{
            "Content-Type":"application/json",
            "Authorization":"Bearer " + token
        },
        body: JSON.stringify({nickname, avatar_url: avatarUrl, bio})
    })
    return await res.json()
}

async function getProfileById(id){
    const res = await fetch(API_BASE + "/profile/" + encodeURIComponent(id))
    return await res.json()
}

async function getVideoById(id){
    const res = await fetch(API_BASE + "/video/" + encodeURIComponent(id))
    return await res.json()
}

async function getMyVideos(){
    const token = getAuthToken()
    if (!token){
        return {error:"login required"}
    }
    const res = await fetch(API_BASE + "/me/videos",{
        headers:{
            "Authorization":"Bearer " + token
        }
    })
    return await res.json()
}

async function getUserVideos(userId){
    const res = await fetch(API_BASE + "/user/" + encodeURIComponent(userId) + "/videos")
    return await res.json()
}

async function updateVideoMeta(videoId, title, description, tags){
    const token = getAuthToken()
    if (!token){
        return {error:"login required"}
    }
    const res = await fetch(API_BASE + "/video/" + encodeURIComponent(videoId),{
        method:"PUT",
        headers:{
            "Content-Type":"application/json",
            "Authorization":"Bearer " + token
        },
        body: JSON.stringify({title, description, tags})
    })
    return await res.json()
}

async function deleteMyVideo(videoId){
    const token = getAuthToken()
    if (!token){
        return {error:"login required"}
    }
    const res = await fetch(API_BASE + "/video/" + encodeURIComponent(videoId),{
        method:"DELETE",
        headers:{
            "Authorization":"Bearer " + token
        }
    })
    return await res.json()
}

async function getMyLikes(){
    const token = getAuthToken()
    if (!token){
        return {error:"login required"}
    }
    const base = getRecommendBase()
    const res = await fetch(base + "/me/likes",{
        headers:{
            "Authorization":"Bearer " + token
        }
    })
    return await res.json()
}

async function getMyFollows(){
    const token = getAuthToken()
    if (!token){
        return {error:"login required"}
    }
    const base = getRecommendBase()
    const res = await fetch(base + "/me/follows",{
        headers:{
            "Authorization":"Bearer " + token
        }
    })
    return await res.json()
}

async function getMyFollowers(){
    const token = getAuthToken()
    if (!token){
        return {error:"login required"}
    }
    const base = getRecommendBase()
    const res = await fetch(base + "/me/followers",{
        headers:{
            "Authorization":"Bearer " + token
        }
    })
    return await res.json()
}

async function getMyHistory(){
    const token = getAuthToken()
    if (!token){
        return {error:"login required"}
    }
    const base = getRecommendBase()
    const res = await fetch(base + "/me/history",{
        headers:{
            "Authorization":"Bearer " + token
        }
    })
    return await res.json()
}

async function getNotifications(){
    const token = getAuthToken()
    if (!token){
        return {error:"login required"}
    }
    const base = getRecommendBase()
    const res = await fetch(base + "/notifications",{
        headers:{
            "Authorization":"Bearer " + token
        }
    })
    return await res.json()
}

async function markNotificationsRead(id, all){
    const token = getAuthToken()
    if (!token){
        return {error:"login required"}
    }
    const base = getRecommendBase()
    const res = await fetch(base + "/notifications/read",{
        method:"POST",
        headers:{
            "Content-Type":"application/json",
            "Authorization":"Bearer " + token
        },
        body: JSON.stringify({id, all: !!all})
    })
    return await res.json()
}

async function addFavorite(videoId){
    const token = getAuthToken()
    if (!token){
        return {error:"login required"}
    }
    const base = getRecommendBase()
    const res = await fetch(base + "/favorite",{
        method:"POST",
        headers:{
            "Content-Type":"application/json",
            "Authorization":"Bearer " + token
        },
        body: JSON.stringify({video_id: videoId})
    })
    return await res.json()
}

async function removeFavorite(videoId){
    const token = getAuthToken()
    if (!token){
        return {error:"login required"}
    }
    const base = getRecommendBase()
    const res = await fetch(base + "/unfavorite",{
        method:"POST",
        headers:{
            "Content-Type":"application/json",
            "Authorization":"Bearer " + token
        },
        body: JSON.stringify({video_id: videoId})
    })
    return await res.json()
}

async function getMyFavorites(){
    const token = getAuthToken()
    if (!token){
        return {error:"login required"}
    }
    const base = getRecommendBase()
    const res = await fetch(base + "/me/favorites",{
        headers:{
            "Authorization":"Bearer " + token
        }
    })
    return await res.json()
}

async function getVideoStats(platformId, videoId){
    const base = getRecommendBase()
    const res = await fetch(base + "/video/" + encodeURIComponent(videoId) + "/stats")
    return await res.json()
}

async function getCommentCount(platformId, videoId){
    const base = getPlatformBase(platformId)
    const res = await fetch(base + "/video/" + encodeURIComponent(videoId) + "/comments/count")
    return await res.json()
}

async function getVideoComments(platformId, videoId){
    const base = getPlatformBase(platformId)
    const res = await fetch(base + "/video/" + videoId + "/comments")
    return await res.json()
}

async function createComment(platformId, videoId, content, parentId){
    const token = getAuthToken()
    if (!token){
        return {error:"login required"}
    }
    const base = getPlatformBase(platformId)
    const res = await fetch(base + "/comment",{
        method:"POST",
        headers:{
            "Content-Type":"application/json",
            "Authorization":"Bearer " + token
        },
        body: JSON.stringify({video_id: videoId, content, parent_id: parentId || 0})
    })
    return await res.json()
}

async function deleteComment(platformId, id){
    const token = getAuthToken()
    if (!token){
        return {error:"login required"}
    }
    const base = getPlatformBase(platformId)
    const res = await fetch(base + "/comment/" + id,{
        method:"DELETE",
        headers:{
            "Authorization":"Bearer " + token
        }
    })
    return await res.json()
}

async function likeComment(platformId, id){
    const token = getAuthToken()
    if (!token){
        return {error:"login required"}
    }
    const base = getPlatformBase(platformId)
    const res = await fetch(base + "/comment/" + id + "/like",{
        method:"POST",
        headers:{
            "Authorization":"Bearer " + token
        }
    })
    return await res.json()
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

async function getVideos(page, limit){
    const params = []
    if (page){
        params.push("page=" + encodeURIComponent(page))
    }
    if (limit){
        params.push("limit=" + encodeURIComponent(limit))
    }
    const url = API_BASE + "/videos" + (params.length ? ("?" + params.join("&")) : "")
    const res = await fetch(url)
    const data = await res.json()
    return data.videos || []
}

function getAuthToken(){

    return localStorage.getItem("token")

}

function logoutUser(){
    localStorage.removeItem("token")
    localStorage.removeItem("user_id")
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

async function uploadVideo(title,description,tags,file,cover){

    const publicKey = localStorage.getItem("public_key")
    const privateKey = localStorage.getItem("private_key")
    if (!publicKey || !privateKey){
        return {error:"keys required, please register first"}
    }
    const videoHash = await hashFile(file)
    const timestamp = Math.floor(Date.now() / 1000)
    const signature = await signMessage(buildProofMessage(videoHash, timestamp, publicKey), privateKey)

    const storageRes = await uploadToStorage(file, videoHash, timestamp, signature, publicKey)
    if (storageRes.error){
        return storageRes
    }

    const storageBase = getStorageBase().replace(/\/+$/,"")
    const manifestUrl = storageRes.manifest_url
        ? (storageRes.manifest_url.startsWith("http") ? storageRes.manifest_url : storageBase + storageRes.manifest_url)
        : ""

    const form = new FormData()
    form.append("title",title)
    form.append("description",description)
    form.append("tags",tags)
    if (cover){
        form.append("cover", cover)
    }
    form.append("storage_id", storageRes.storage_id || "")
    form.append("manifest_url", manifestUrl)
    form.append("manifest_hash", storageRes.manifest_hash || "")
    form.append("chunks", JSON.stringify(storageRes.chunks || []))
    form.append("video_hash", videoHash)
    form.append("author_timestamp", String(timestamp))
    form.append("author_signature", signature)
    if (storageRes.filename){
        form.append("filename", storageRes.filename)
    } else if (file && file.name) {
        form.append("filename", file.name)
    }

    const token = getAuthToken()

    const headers = {}
    if (token){
        headers["Authorization"] = "Bearer " + token
    }

    const res = await fetch(API_BASE + "/upload",{
        method:"POST",
        body:form,
        headers
    })

    return await res.json()

}

async function uploadToStorage(file, videoHash, timestamp, signature, publicKey){
    const base = getStorageBase()
    const form = new FormData()
    form.append("file", file)
    form.append("video_hash", videoHash)
    form.append("author_timestamp", String(timestamp))
    form.append("author_signature", signature)
    form.append("author_public_key", publicKey)
    const res = await fetch(base + "/store",{
        method:"POST",
        body: form
    })
    return await res.json()
}

async function getRecommend(user, recType, page, limit){

    const base = getRecommendBase()
    const params = []
    if (user){
        params.push("user=" + encodeURIComponent(user))
    }
    if (recType){
        params.push("type=" + encodeURIComponent(recType))
    }
    if (page){
        params.push("page=" + encodeURIComponent(page))
    }
    if (limit){
        params.push("limit=" + encodeURIComponent(limit))
    }
    const url = base + "/recommend" + (params.length ? ("?" + params.join("&")) : "")

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

    const base = getRecommendBase()
    const res = await fetch(base + "/behavior",{
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

async function searchVideos(q, page, limit){
    const params = ["q=" + encodeURIComponent(q)]
    if (page){
        params.push("page=" + encodeURIComponent(page))
    }
    if (limit){
        params.push("limit=" + encodeURIComponent(limit))
    }
    const base = API_BASE
    const res = await fetch(base + "/search?" + params.join("&"))
    const data = await res.json()
    return data.videos || []
}

async function searchVideosByTag(tag, page, limit){
    const params = ["tag=" + encodeURIComponent(tag)]
    if (page){
        params.push("page=" + encodeURIComponent(page))
    }
    if (limit){
        params.push("limit=" + encodeURIComponent(limit))
    }
    const base = API_BASE
    const res = await fetch(base + "/search?" + params.join("&"))
    return await res.json()
}

async function searchUsers(q){
    const res = await fetch(API_BASE + "/users/search?q=" + encodeURIComponent(q))
    return await res.json()
}

async function registerUser(publicKey){

    const res = await fetch(API_BASE + "/register",{
        method:"POST",
        headers:{"Content-Type":"application/json"},
        body: JSON.stringify({public_key: publicKey})
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

function bytesToHex(bytes){
    const hex = []
    for (let i=0;i<bytes.length;i++){
        hex.push(bytes[i].toString(16).padStart(2,"0"))
    }
    return hex.join("")
}

async function generateKeypair(){
    const keypair = await crypto.subtle.generateKey(
        {name:"Ed25519"},
        true,
        ["sign","verify"]
    )
    const priv = await crypto.subtle.exportKey("pkcs8", keypair.privateKey)
    const pub = await crypto.subtle.exportKey("spki", keypair.publicKey)
    return {
        public_key: bytesToB64(new Uint8Array(pub)),
        private_key: bytesToB64(new Uint8Array(priv))
    }
}

function readKeyBundle(){
    const publicKey = localStorage.getItem("public_key")
    const privateKey = localStorage.getItem("private_key")
    if (!publicKey || !privateKey){
        return null
    }
    return {public_key: publicKey, private_key: privateKey}
}

function exportKeyBundle(){
    const bundle = readKeyBundle()
    if (!bundle){
        return null
    }
    return JSON.stringify(bundle)
}

function importKeyBundle(jsonText){
    if (!jsonText || !jsonText.trim()){
        return {error:"key json required"}
    }
    let data
    try {
        data = JSON.parse(jsonText)
    } catch (e) {
        return {error:"invalid json"}
    }
    if (!data || typeof data.public_key !== "string" || typeof data.private_key !== "string"){
        return {error:"missing public_key/private_key"}
    }
    localStorage.setItem("public_key", data.public_key)
    localStorage.setItem("private_key", data.private_key)
    return {ok:true}
}

function downloadKeyBundle(){
    const payload = exportKeyBundle()
    if (!payload){
        return {error:"no keys in local storage"}
    }
    const blob = new Blob([payload], {type:"application/json"})
    const url = URL.createObjectURL(blob)
    const a = document.createElement("a")
    a.href = url
    a.download = "keys.json"
    document.body.appendChild(a)
    a.click()
    a.remove()
    URL.revokeObjectURL(url)
    return {ok:true}
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

async function signMessage(message, privateKeyB64){
    const keyBytes = b64ToBytes(privateKeyB64)
    const key = await crypto.subtle.importKey(
        "pkcs8",
        keyBytes,
        {name:"Ed25519"},
        false,
        ["sign"]
    )
    const msgBytes = new TextEncoder().encode(message)
    const sig = await crypto.subtle.sign("Ed25519", key, msgBytes)
    return bytesToB64(new Uint8Array(sig))
}

async function hashFile(file){
    const buf = await file.arrayBuffer()
    const hash = await crypto.subtle.digest("SHA-256", buf)
    return bytesToHex(new Uint8Array(hash))
}

function buildProofMessage(videoHash, timestamp, publicKey){
    return videoHash + "|" + timestamp + "|" + publicKey
}

async function loginUser(publicKey){

    const challengeRes = await fetch(API_BASE + "/login/challenge",{
        method:"POST",
        headers:{"Content-Type":"application/json"},
        body: JSON.stringify({public_key: publicKey})
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
        body: JSON.stringify({public_key: publicKey, nonce, signature})
    })

    return await res.json()

}

async function followAuthor(authorId){

    const token = getAuthToken()
    if (!token){
        return {error:"login required"}
    }

    const base = getRecommendBase()
    const res = await fetch(base + "/follow",{
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

    const base = getRecommendBase()
    const res = await fetch(base + "/unfollow",{
        method:"POST",
        headers:{
            "Content-Type":"application/json",
            "Authorization":"Bearer " + token
        },
        body: JSON.stringify({author_id: authorId})
    })

    return await res.json()

}
