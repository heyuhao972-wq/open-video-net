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

const streamBase = getPlatformBase(platform)
player.src = streamBase + "/video/" + encodeURIComponent(id) + "/stream"
player.preload = "auto"

if (id) {
    reportBehavior(id,"watch")
}

const resumeKey = "resume_" + platform + "_" + id
player.addEventListener("loadedmetadata", ()=>{
    const saved = localStorage.getItem(resumeKey)
    const seconds = saved ? parseFloat(saved) : 0
    if (seconds && !isNaN(seconds) && seconds < player.duration - 5) {
        player.currentTime = seconds
    }
})

let lastSave = 0
player.addEventListener("timeupdate", ()=>{
    const now = Date.now()
    if (now - lastSave < 5000) {
        return
    }
    lastSave = now
    localStorage.setItem(resumeKey, String(player.currentTime))
})

player.addEventListener("ended", ()=>{
    localStorage.removeItem(resumeKey)
})

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

async function loadStats(){
    const box = document.getElementById("stats")
    if (!box || !id){
        return
    }
    const stats = await getVideoStats(platform, id)
    const commentCount = await getCommentCount(platform, id)
    if (stats.error){
        box.innerText = stats.error
        return
    }
    const comments = commentCount && commentCount.count ? commentCount.count : 0
    box.innerText = "views: " + (stats.watches || 0) + " | likes: " + (stats.likes || 0) + " | shares: " + (stats.shares || 0) + " | favorites: " + (stats.favorites || 0) + " | comments: " + comments
}

loadStats()

async function loadComments(){
    const list = document.getElementById("comment-list")
    if (!list || !id){
        return
    }
    const data = await getVideoComments(platform, id)
    list.innerHTML = ""
    if (!Array.isArray(data)){
        return
    }
    data.forEach(c=>{
        const li = document.createElement("li")
        const text = document.createElement("span")
        text.innerText = (c.user_id || "user") + ": " + (c.content || "") + " (" + (c.likes || 0) + ")"
        li.appendChild(text)

        const likeBtn = document.createElement("button")
        likeBtn.innerText = "Like"
        likeBtn.addEventListener("click", async ()=>{
            const res = await likeComment(platform, c.id)
            if (res.error){
                alert(res.error)
                return
            }
            await loadComments()
            await loadStats()
        })
        li.appendChild(likeBtn)

        const delBtn = document.createElement("button")
        delBtn.innerText = "Delete"
        delBtn.addEventListener("click", async ()=>{
            const res = await deleteComment(platform, c.id)
            if (res.error){
                alert(res.error)
                return
            }
            await loadComments()
            await loadStats()
        })
        li.appendChild(delBtn)

        list.appendChild(li)
    })
}

const commentForm = document.getElementById("comment-form")
if (commentForm){
    commentForm.addEventListener("submit", async e=>{
        e.preventDefault()
        const input = document.getElementById("comment-input")
        const content = (input && input.value) ? input.value.trim() : ""
        if (!content){
            return
        }
        const res = await createComment(platform, id, content)
        if (res.error){
            alert(res.error)
            return
        }
        reportBehavior(id,"comment")
        if (input){
            input.value = ""
        }
        await loadComments()
        await loadStats()
    })
}

loadComments()

const likeBtn = document.getElementById("like-btn")
if (likeBtn){
    likeBtn.addEventListener("click", ()=>{
        if (id){
            reportBehavior(id,"like")
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
    })
}

const niBtn = document.getElementById("not-interested-btn")
if (niBtn){
    niBtn.addEventListener("click", ()=>{
        if (id){
            reportBehavior(id,"not_interested")
        }
    })
}

const favBtn = document.getElementById("favorite-btn")
if (favBtn){
    favBtn.addEventListener("click", async ()=>{
        if (!id){
            return
        }
        const res = await addFavorite(id)
        if (res.error){
            alert(res.error)
            return
        }
        await loadStats()
    })
}

const unfavBtn = document.getElementById("unfavorite-btn")
if (unfavBtn){
    unfavBtn.addEventListener("click", async ()=>{
        if (!id){
            return
        }
        const res = await removeFavorite(id)
        if (res.error){
            alert(res.error)
            return
        }
        await loadStats()
    })
}

const followBtn = document.getElementById("follow-btn")
if (followBtn){
    followBtn.addEventListener("click", async ()=>{
        const meta = document.getElementById("meta")
        const authorId = meta.dataset.authorId
        if (!authorId){
            alert("author id missing")
            return
        }
        const res = await followAuthor(authorId)
        if (res.error){
            alert(res.error)
            return
        }
        alert("followed")
    })
}

const unfollowBtn = document.getElementById("unfollow-btn")
if (unfollowBtn){
    unfollowBtn.addEventListener("click", async ()=>{
        const meta = document.getElementById("meta")
        const authorId = meta.dataset.authorId
        if (!authorId){
            alert("author id missing")
            return
        }
        const res = await unfollowAuthor(authorId)
        if (res.error){
            alert(res.error)
            return
        }
        alert("unfollowed")
    })
}
