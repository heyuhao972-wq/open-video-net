const form = document.getElementById("upload-form")

form.addEventListener("submit",async e=>{

    e.preventDefault()

    const title = document.getElementById("title").value

    const description = document.getElementById("description").value

    const tags = document.getElementById("tags").value

    const cover = document.getElementById("cover").files[0]
    const file = document.getElementById("file").files[0]

    const res = await uploadVideo(title,description,tags,file,cover)

    if (res && res.error){
        alert(res.error)
        return
    }
    alert("upload success")

})
