export default class PostController {
  formEl;
  boxEl;
  constructor(form, box) {
    this.addEvent(form, box);
  }

  addEvent(form, box) {
    this.formEl = document.querySelector(form);
    if (!this.formEl) {
      console.error(`${form} element not found`);
      return;
    }
    this.boxEl = this.formEl.querySelector(box);
    this.formEl.addEventListener("submit", (e) => {
      e.preventDefault();
      this.creatPost();
    });
  }
  async creatPost() {
    const formData = new FormData(this.formEl);
    let tags = tagsToArray(formData.get("tags"));
    const data = {
      title: formData.get("title"),
      content: formData.get("content"),
      tags: tags,
    };
    if (tags.length == 0 || data.title == "" || data.content == "") {
      this.writeError("All fields must be completed.");
      return;
    }
    try {
      const response = await fetch("/api/post", {
        method: "POST",
        body: JSON.stringify(data),
      });
      let responseData = await response.json();
      if (response.ok) {
        this.writeSucc(responseData.message);
        setTimeout(() => {
          window.location.href = "/";
        }, 1000);
      } else {
        this.writeError(responseData.message);
      }
    } catch (error) {
      console.log(error);
    }
  }
  writeError(msg) {
    this.boxEl.classList.remove("success");
    this.boxEl.classList.add("error");
    this.boxEl.innerHTML = msg;
  }

  writeSucc(msg) {
    this.boxEl.classList.remove("error");
    this.boxEl.classList.add("success");
    this.boxEl.innerHTML = msg;
  }
}

function tagsToArray(categoryString) {
  const categoriesArray = categoryString
    .split(",")
    .map((category) => category.trim())
    .filter((cate) => cate.length > 0);
  return categoriesArray;
}
