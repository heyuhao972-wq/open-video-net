async function loadProfile(){
    const box = document.getElementById("profile-info")
    if (!box){
        return
    }
    const res = await getMyProfile()
    if (res.error){
        box.innerText = res.error
        return
    }
    const u = res.user || {}
    box.innerText = "id: " + (u.id || "") + "\npublic_key: " + (u.public_key || "") + "\nnickname: " + (u.nickname || "") + "\navatar: " + (u.avatar_url || "") + "\nbio: " + (u.bio || "")
}

async function loadMyVideos(){
    const list = document.getElementById("my-videos")
    if (!list){
        return
    }
    const res = await getMyVideos()
    if (res.error){
        list.innerText = res.error
        return
    }
    list.innerHTML = ""
    const videos = res.videos || []
    videos.forEach(v=>{
        const li = document.createElement("li")
        const link = document.createElement("a")
        const platformId = v.platform_id || "platformA"
        link.href = "player.html?platform=" + encodeURIComponent(platformId) + "&id=" + v.id
        link.innerText = v.title || v.id
        li.appendChild(link)

        const editBtn = document.createElement("button")
        editBtn.innerText = "Edit"
        editBtn.addEventListener("click", ()=>{
            document.getElementById("edit-id").value = v.id
            document.getElementById("edit-title").value = v.title || ""
            document.getElementById("edit-desc").value = v.description || ""
            document.getElementById("edit-tags").value = (v.tags || []).join(",")
        })
        li.appendChild(editBtn)

        const delBtn = document.createElement("button")
        delBtn.innerText = "Delete"
        delBtn.addEventListener("click", async ()=>{
            const r = await deleteMyVideo(v.id)
            if (r.error){
                alert(r.error)
                return
            }
            await loadMyVideos()
        })
        li.appendChild(delBtn)

        list.appendChild(li)
    })
}

async function loadMyLikes(){
    const list = document.getElementById("my-likes")
    if (!list){
        return
    }
    const res = await getMyLikes()
    if (res.error){
        list.innerText = res.error
        return
    }
    list.innerHTML = ""
    const videos = res.videos || []
    videos.forEach(id=>{
        const li = document.createElement("li")
        li.innerText = id
        list.appendChild(li)
    })
}

async function loadMyFavorites(){
    const list = document.getElementById("my-favorites")
    if (!list){
        return
    }
    const res = await getMyFavorites()
    if (res.error){
        list.innerText = res.error
        return
    }
    list.innerHTML = ""
    const videos = res.videos || []
    videos.forEach(id=>{
        const li = document.createElement("li")
        li.innerText = id
        list.appendChild(li)
    })
}

async function loadMyHistory(){
    const list = document.getElementById("my-history")
    if (!list){
        return
    }
    const res = await getMyHistory()
    if (res.error){
        list.innerText = res.error
        return
    }
    list.innerHTML = ""
    const items = res.history || []
    items.forEach(h=>{
        const li = document.createElement("li")
        const platformId = h.platform_id || "platformA"
        const link = document.createElement("a")
        link.href = "player.html?platform=" + encodeURIComponent(platformId) + "&id=" + h.video_id
        link.innerText = h.video_id
        li.appendChild(link)
        list.appendChild(li)
    })
}

async function loadMyFollows(){
    const list = document.getElementById("my-follows")
    if (!list){
        return
    }
    const res = await getMyFollows()
    if (res.error){
        list.innerText = res.error
        return
    }
    list.innerHTML = ""
    const users = res.users || []
    users.forEach(id=>{
        const li = document.createElement("li")
        li.innerText = id
        list.appendChild(li)
    })
}

async function loadMyFollowers(){
    const list = document.getElementById("my-followers")
    if (!list){
        return
    }
    const res = await getMyFollowers()
    if (res.error){
        list.innerText = res.error
        return
    }
    list.innerHTML = ""
    const users = res.users || []
    users.forEach(id=>{
        const li = document.createElement("li")
        li.innerText = id
        list.appendChild(li)
    })
}

const form = document.getElementById("nickname-form")
if (form){
    form.addEventListener("submit", async e=>{
        e.preventDefault()
        const nick = document.getElementById("nickname-input").value
        const avatar = document.getElementById("avatar-input").value
        const bio = document.getElementById("bio-input").value
        const res = await updateProfile(nick, avatar, bio)
        if (res.error){
            alert(res.error)
            return
        }
        await loadProfile()
    })
}

loadProfile()
loadMyVideos()
loadMyLikes()
loadMyFavorites()
loadMyHistory()
loadMyFollows()
loadMyFollowers()

const saveBtn = document.getElementById("edit-save")
if (saveBtn){
    saveBtn.addEventListener("click", async ()=>{
        const id = document.getElementById("edit-id").value
        const title = document.getElementById("edit-title").value
        const desc = document.getElementById("edit-desc").value
        const tags = document.getElementById("edit-tags").value
        if (!id){
            alert("video id required")
            return
        }
        const res = await updateVideoMeta(id, title, desc, tags)
        if (res.error){
            alert(res.error)
            return
        }
        await loadMyVideos()
    })
}

const logoutBtn = document.getElementById("logout-btn")
if (logoutBtn){
    logoutBtn.addEventListener("click", e=>{
        e.preventDefault()
        logoutUser()
        window.location.href = "login.html"
    })
}
