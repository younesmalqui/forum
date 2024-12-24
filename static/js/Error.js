export class DOMError {
  errorEL;
  constructor(el) {
    this.el = el;
    this.errorEL = this.createElement();
  }

  createElement() {
    let div = document.createElement("div");
    this.el.prepend(div);
    return div;
  }

  writeError(msg) {
    this.errorEL.classList.remove("success");
    this.errorEL.classList.add("error");
    this.errorEL.innerHTML = msg;
  }

  writeSucc(msg) {
    this.errorEL.classList.remove("error");
    this.errorEL.classList.add("success");
    this.errorEL.innerHTML = msg;
  }
  
  remove() {
    this.el.remove();
  }
}
