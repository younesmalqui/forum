import { likePost, dislikePost, getPostReactions } from "./api.js";
import { redirectTo } from "./utils.js";
export default class PostReact {
  constructor(parentEl) {
    /** @type {HTMLDivElement} */
    this.parentEl = parentEl;
    this.initEvents();
  }

  initEvents() {
    this.parentEl.addEventListener("click", this.handleOnClickEvent.bind(this));
  }

  /**
   *
   * @param {PointerEvent} e
   */
  handleOnClickEvent(e) {
    const postElement = e.target.closest(".post");
    if (!postElement) return;
    const postView = new PostView(postElement, e.target);
    if (postView.isDislikeClicked) {
      e.preventDefault();
      return this.handleDislike(postView);
    }
    if (postView.isLikeClicked) {
      e.preventDefault();
      return this.handleLike(postView);
    }
  }
  /**
   *
   * @param {PostView} postView
   */
  async handleDislike(postView) {
    try {
      await dislikePost(postView.postId);
      const postData = await getPostReactions(postView.postId);
      const { likesCount, dislikesCount } = postData.data;
      postView.updateReactions(likesCount, dislikesCount);
    } catch (error) {
      if (error?.status == 401) {
        redirectTo("/login");
        return;
      }
      console.error(error);
      postView.updateDisLikeError();
    }
  }

  async handleLike(postView) {
    try {
      await likePost(postView.postId);
      const postData = await getPostReactions(postView.postId);
      const { likesCount, dislikesCount } = postData.data;
      postView.updateReactions(likesCount, dislikesCount);
    } catch (error) {
      console.error(error);
      postView.updateDisLikeError();
    }
  }
}

class PostView {
  /**
   *
   * @param {HTMLDivElement} postElement
   * @param {HTMLElement} target
   */
  constructor(postElement, target) {
    this.postElement = postElement;
    this.target = target;
    this.initElements();
  }

  initElements() {
    this.likeUp = this.postElement.querySelector(".like-up");
    this.likeDown = this.postElement.querySelector(".like-down");
    this.likeCount = this.likeUp.querySelector(".like-count");
    this.dislikeCount = this.likeDown.querySelector(".like-count");
  }

  get postId() {
    return +this.postElement.dataset.id;
  }

  get isDislikeClicked() {
    return this.target.closest(".like-down");
  }

  get isLikeClicked() {
    return this.target.closest(".like-up");
  }

  updateReactions(likes, dislikes) {
    this.likeCount.innerHTML = likes;
    this.dislikeCount.innerHTML = dislikes;
  }

  updateLikeError() {
    this.likeCount.style.color = "red";
    setTimeout(() => {
      this.likeCount.style.color = "";
    }, 1000);
  }

  updateDisLikeError() {
    this.dislikeCount.style.color = "red";
    setTimeout(() => {
      console.log("done");
      this.dislikeCount.style.color = "";
    }, 1000);
  }
}
