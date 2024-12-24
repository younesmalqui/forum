export async function likePost(postId) {
  const bodyData = {
    postId,
    isLike: 1,
    action: "react",
  };
  const response = await fetch(`/api/react`, {
    method: "POST",
    body: JSON.stringify(bodyData),
  });
  const data = await response.json();
  if (!response.ok) throw data;
  return data;
}

export async function dislikePost(postId) {
  const bodyData = {
    postId,
    isLike: -1,
    action: "react",
  };
  const response = await fetch(`/api/react`, {
    method: "POST",
    body: JSON.stringify(bodyData),
  });
  const data = await response.json();
  if (!response.ok) throw data;
  return data;
}

export async function getPostReactions(postId) {
  const response = await fetch(`/api/react?postId=${postId}`);
  const data = await response.json();
  if (!response.ok) throw data;
  return data;
}

//------------- COMMENTS

export async function addComment(body) {
  const response = await fetch("/api/add/comment", {
    method: "POST",
    body: JSON.stringify(body),
  });

  const data = await response.json();
  if (!response.ok) {
    throw data;
  }
  return data;
}

export async function reactComment(commentId, isLike) {
  const body = {
    isLike,
    commentId,
  };
  const response = await fetch("/api/like/comment", {
    method: "POST",
    body: JSON.stringify(body),
  });
  const data = await response.json();
  if (!response.ok) {
    throw data;
  }
  return data;
}
