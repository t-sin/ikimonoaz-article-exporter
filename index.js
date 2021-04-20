const axios = require('axios');

const api_base_url = 'https://ikimonoaz.ikimonopal.jp/api'

// General purpose GET function to IkimonoAZ API.
const get_from_api = async (path) => {
  const url = `${api_base_url}/${path}`;
  const ret = await axios.get(url);

  if (ret.request.res.statusCode !== 200) {
    throw new Error(`server says error: ${ret.request.res.statusCode}`);
  };

  return ret.data;
}

// Simplify article data.
const map_to_simple_article = (raw_article) => {
  return {
    // 記事ID
    article_id: raw_article.article_id,
    // 日時
    created_at: raw_article.create_date,
    updated_at: raw_article.edit_date,
    released_at: raw_article.release_date,
    // 記事データ
    title: raw_article.title,
    contents: raw_article.contents,
    // 写真データ
    media: raw_article.media.map((e) => ({ url: e.url, movie_thumbnail: e.movie_thubnail })),
    // メタデータ
    creatures: raw_article.creatures.map((e) => `${e.name} @ ${e.place.name}`),
    tags: raw_article.tags.map((e) => e.name),
  };
};

const map_to_simple_user = (raw_user) => {
  return {
    // ユーザ情報
    name: raw_user.user_name,
    profile: raw_user.profile,
    profile_image: raw_user.profile_image_url,
    // マイスター一覧
    meister: raw_user.meister.map((e) => e.name),
    // よく行く園館
    place_name: raw_user.place_name,
  };
};

// Get one article data.
const get_article = async (article_id) => {
  const path = `articles/get?articleId=${article_id}&incView=1`
  console.log(`retrieving article data from '${path}`);
  const body = await get_from_api(path);

  return map_to_simple_article(body.data);
};

// Get pagenated articles.
const collect_articles = async (user_id, page) => {
  const path = `articles/list?mypageUserId=${user_id}&page=${page}&stripHtml=0`;
  console.log(`collecting articles from '${path}`);
  const body = await get_from_api(path);

  return body.data.articles.map(map_to_simple_article);
};

// Get comments in specified article ID.
const collect_comments = async (article_id) => {
  const path = `comments/list?articleId=${article_id}&offset=0&isFuture=false`
  console.log(`collecting comments from '${path}`);
  const body = await get_from_api(path);

  return body.data.comments;
};

// Get user info by user ID.
const get_user_info = async (user_id) => {
  const path = `users/get?userId=${user_id}`
  console.log(`get user info from '${path}`);
  const body = await get_from_api(path);

  return map_to_simple_user(body.data);
};

// Get all articles by user ID.
const collect_all_articles = async (user_id) => {
  let page = 1;
  let articles = [];
  let result;

  do {
    console.log(`page ${page}...`);

    result = await collect_articles(user_id, page);
    articles = articles.concat(result);
    page += 1;
  } while (result.length != 0);

  return articles;
}

const collect_all_user_data = async (user_id) => {
  const user = await get_user_info(user_id);
  let articles = await collect_all_articles(user_id);

  for (let a of articles) {
    const comments = await collect_comments(a.article_id);
    const simple_comments = comments
      .filter((c) => c.user_id !== null)
      .map((c) => {
        console.log(c);
        return { comment: c.comment, user_name: c.user.user_name };
      });
    a.comments = simple_comments;
  }
  
  return {
    user: user,
    articles: articles,
  };
};

///////////////////////////////////////////////
const user_id = 7308; // 自分のID

collect_all_user_data(user_id).then((articles) => {
  console.log(articles);
}).catch((e) => {
  console.log(e);
});
