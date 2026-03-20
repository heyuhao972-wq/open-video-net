const followForm = document.getElementById("follow-form")

followForm.addEventListener("submit", async e=>{

    e.preventDefault()

    const authorId = document.getElementById("author-id").value
    const res = await followAuthor(authorId)
    if (res.error){
        alert(res.error)
        return
    }
    alert("followed")
})
