const form = document.getElementById("upload-form")

form.addEventListener("submit",async e=>{

    e.preventDefault()

    const status = document.getElementById("upload-status")
    if (status){
        status.classList.remove("error")
        status.innerText = "Uploading..."
    }

    const title = document.getElementById("title").value

    const description = document.getElementById("description").value

    const tags = document.getElementById("tags").value

    const platformBase = document.getElementById("platform-select").value
    setContentBase(platformBase)

    const file = document.getElementById("file").files[0]

    const privateKey = localStorage.getItem("private_key")
    const publicKey = localStorage.getItem("public_key")
    if (!privateKey || !publicKey){
        if (status){
            status.classList.add("error")
            status.innerText = "missing keypair, please register"
        }
        return
    }

    const videoHash = await hashFile(file)
    const timestamp = Math.floor(Date.now() / 1000)
    const msg = buildProofMessage(videoHash, timestamp, publicKey)
    const signature = await signMessage(msg, privateKey)

    const res = await uploadVideo(title,description,tags,file,{
        signature,
        timestamp,
        videoHash
    })
    if (res && res.error){
        if (status){
            status.classList.add("error")
            status.innerText = res.error
        } else {
            alert(res.error)
        }
        return
    }

    if (status){
        status.innerText = "Upload success"
    } else {
        alert("upload success")
    }

})
