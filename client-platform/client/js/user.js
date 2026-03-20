function getQuery(name){
    const params = new URLSearchParams(window.location.search)
    return params.get(name)
}

async function loadUser(){
    const id = getQuery("id")
    const info = document.getElementById("user-info")
    const list = document.getElementById("user-videos")
    if (!id || !info || !list){
        return
    }
    const res = await getProfileById(id)
    if (res.error){
        info.innerText = res.error
        return
    }
    const u = res.user || {}
    info.innerText = "nickname: " + (u.nickname || "") + "\nuser_id: " + (u.id || "") + "\navatar: " + (u.avatar_url || "") + "\nbio: " + (u.bio || "")

    const vids = await getUserVideos(id)
    const videos = vids.videos || []
    list.innerHTML = ""
    videos.forEach(v=>{
        const li = document.createElement("li")
        const link = document.createElement("a")
        const platformId = v.platform_id || "platformA"
        link.href = "player.html?platform=" + encodeURIComponent(platformId) + "&id=" + v.id
        link.innerText = v.title || v.id
        li.appendChild(link)
        list.appendChild(li)
    })
}

loadUser()
