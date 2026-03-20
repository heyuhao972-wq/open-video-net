const loginForm = document.getElementById("login-form")
const exportBtn = document.getElementById("export-keys")
const importInput = document.getElementById("import-keys")
const publicKeyBox = document.getElementById("public-key")

loginForm.addEventListener("submit", async e=>{

    e.preventDefault()

    const publicKey = localStorage.getItem("public_key")
    if (!publicKey){
        alert("public key missing, please register")
        return
    }
    if (publicKeyBox){
        publicKeyBox.value = publicKey
    }
    const res = await loginUser(publicKey)

    if (res.error){
        alert(res.error)
        return
    }

    localStorage.setItem("token", res.token)
    localStorage.setItem("user_id", res.user.id)

    window.location.href = "index.html"

})

if (exportBtn){
    exportBtn.addEventListener("click", ()=>{
        const publicKey = localStorage.getItem("public_key") || ""
        const privateKey = localStorage.getItem("private_key") || ""
        if (!publicKey || !privateKey){
            alert("no keys found")
            return
        }
        const payload = JSON.stringify({public_key: publicKey, private_key: privateKey}, null, 2)
        const blob = new Blob([payload], {type: "application/json"})
        const url = URL.createObjectURL(blob)
        const a = document.createElement("a")
        a.href = url
        a.download = "openvideo-keys.json"
        a.click()
        URL.revokeObjectURL(url)
    })
}

if (importInput){
    importInput.addEventListener("change", async e=>{
        const file = e.target.files[0]
        if (!file){
            return
        }
        const text = await file.text()
        try {
            const data = JSON.parse(text)
            if (!data.public_key || !data.private_key){
                alert("invalid key file")
                return
            }
            localStorage.setItem("public_key", data.public_key)
            localStorage.setItem("private_key", data.private_key)
            if (publicKeyBox){
                publicKeyBox.value = data.public_key
            }
            alert("keys imported")
        } catch (err) {
            alert("invalid key file")
        }
    })
}
