import { addComment, reactComment } from "./api.js";
import { DOMError } from "./Error.js";

export class CommentForm {
  constructor(selector) {
    this.formEl = this.initFormEl(selector);
    this.initEvent();
    this.alertEl = new DOMError(this.formEl);
  }

  initFormEl(selector) {
    const el = document.querySelector(selector);
    if (!el) throw new Error("comment form is not found");
    return el;
  }

  initEvent() {
    this.formEl.addEventListener("submit", this.handleOnSubmit.bind(this));
  }

  /**
   *
   * @param {SubmitEvent} e
   */
  async handleOnSubmit(e) {
    e.preventDefault();
    this.formData = new FormData(this.formEl);
    if (this.content.trim().length == 0) {
      this.alertEl.writeError("Comment cannot be empty");
      return;
    }
    try {
      const data = await addComment(this.data);
      this.alertEl.writeSucc(data.message);
      this.reload();
    } catch (error) {
      this.alertEl.writeError(error.message);
    }
  }

  reload(time = 1000) {
    setTimeout(() => {
      window.location.href = "";
    }, time);
  }

  get content() {
    return this.formData.get("comment");
  }

  get postId() {
    return +this.formEl.dataset.postid;
  }

  get data() {
    return {
      comment: this.content,
      postId: this.postId,
    };
  }
}

class Comment {
  constructor(el) {
    this.el = el;
    this.likeEl = el.querySelector(".like-up").closest(".like-box");
    this.dislikeEl = el.querySelector(".like-down").closest(".like-box");
  }

  updateReaction(like, dislike) {
    console.log(this.likeEl, dislike);
    this.likeEl.querySelector(".like-count").innerText = like;
    this.dislikeEl.querySelector(".like-count").innerText = dislike;
  }

  displayError(error) {
    this.likeEl.querySelector("span").style.color = "red";
    this.dislikeEl.querySelector("span").style.color = "red";
    this.resetError();
    console.error(error);
  }

  resetError(time = 500) {
    setTimeout(() => {
      this.likeEl.querySelector("span").style.color = "";
      this.dislikeEl.querySelector("span").style.color = "";
    }, time);
  }

  get commentId() {
    return +this.el.dataset.id;
  }
}

export class CommentLike {
  constructor(selector) {
    this.selector = selector;
  }

  init() {
    this.commentList = document.querySelector(this.selector);
    if (!this.commentList)
      throw new Error(`element with ${this.selector} selector doesnt exists`);
    this.initEvent();
  }

  initEvent() {
    this.commentList.addEventListener("click", this.handleOnClick.bind(this));
  }

  /**
   *
   * @param {PointerEvent} e
   */
  async handleOnClick(e) {
    e.preventDefault();
    let commentEl = e.target.closest(".comment");
    if (!commentEl) return;
    let comment = new Comment(commentEl);
    let value = this.getValueReaction(e.target);
    if (value === null) return;
    try {
      const { data } = await reactComment(comment.commentId, value);
      comment.updateReaction(data.likes, data.dislikes);
    } catch (error) {
      comment.displayError(error);
    }
  }

  getCommentEl(target) {
    const commentEl = target.closest(".comment");
    return commentEl;
  }

  getValueReaction(target) {
    if (target.closest(".like-up")) return 1;
    if (target.closest(".like-down")) return -1;
    return null;
  }
}

new CommentForm("#commentForm");
const commentLike = new CommentLike("#comment-list");
commentLike.init();
