function getQuery(name){

    const params = new URLSearchParams(window.location.search)

    return params.get(name)

}

let id = getQuery("id")
let platform = getQuery("platform") || "platformA"
const uri = getQuery("uri")
if (uri){
    const info = parseVideoURI(decodeURIComponent(uri))
    if (info){
        id = info.videoId
        platform = info.platformId
    }
}

const player = document.getElementById("player")
const statusBox = document.getElementById("status")

const base = getPlatformBase(platform)
player.src = base + "/video/" + id + "/stream"

if (id) {
    reportBehavior(id,"watch")
}

async function loadMeta(){
    if (!id){
        return
    }
    const base = getPlatformBase(platform)
    const res = await fetch(base + "/video/" + id)
    const data = await res.json()
    const meta = document.getElementById("meta")
    if (data && data.title){
        meta.innerText = "Title: " + data.title + " | Author: " + (data.author_id || "unknown")
        meta.dataset.authorId = data.author_id || ""
    }
}

loadMeta()

const likeBtn = document.getElementById("like-btn")
if (likeBtn){
    likeBtn.addEventListener("click", ()=>{
        if (id){
            reportBehavior(id,"like")
            showStatus("Liked")
        }
    })
}

const shareBtn = document.getElementById("share-btn")
if (shareBtn){
    shareBtn.addEventListener("click", async ()=>{
        if (!id){
            return
        }
        try {
            await navigator.clipboard.writeText(window.location.href)
        } catch (e) {
            // ignore
        }
        reportBehavior(id,"share")
        showStatus("Shared link copied")
    })
}

const niBtn = document.getElementById("not-interested-btn")
if (niBtn){
    niBtn.addEventListener("click", ()=>{
        if (id){
            reportBehavior(id,"not_interested")
            showStatus("Marked as not interested")
        }
    })
}

const followBtn = document.getElementById("follow-btn")
if (followBtn){
    followBtn.addEventListener("click", async ()=>{
        const meta = document.getElementById("meta")
        const authorId = meta.dataset.authorId
        if (!authorId){
            showStatus("Author ID missing", true)
            return
        }
        const res = await followAuthor(authorId)
        if (res.error){
            showStatus(res.error, true)
            return
        }
        showStatus("Followed author")
    })
}

const unfollowBtn = document.getElementById("unfollow-btn")
if (unfollowBtn){
    unfollowBtn.addEventListener("click", async ()=>{
        const meta = document.getElementById("meta")
        const authorId = meta.dataset.authorId
        if (!authorId){
            showStatus("Author ID missing", true)
            return
        }
        const res = await unfollowAuthor(authorId)
        if (res.error){
            showStatus(res.error, true)
            return
        }
        showStatus("Unfollowed author")
    })
}

function showStatus(text, isError){
    if (!statusBox){
        return
    }
    statusBox.classList.toggle("error", !!isError)
    statusBox.innerText = text
}
