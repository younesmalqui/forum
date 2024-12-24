import PostReact from "/static/js/React.js";
import PostController from "/static/js/PostController.js";
let postController = new PostController("#new_post_form", ".box");
let postReact = new PostReact(posts);
