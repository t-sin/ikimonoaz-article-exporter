const output = require('./output.js');
const api = require('./api.js');

const parse5 = require('parse5');
import 'regenerator-runtime/runtime';

const parse_contents = (contents) => {
  return parse5.parseFragment(contents);
};

// Iterate on media nodes and do `fn`.
const for_each_media_nodes = (doc, fn) => {
  const nodes = doc.childNodes;
  const media_nodes = nodes.filter((n) => n.nodeName == 'img' || n.nodeName == 'video');
  media_nodes.forEach(fn);
};

// Replace media link as local-referencible link.
const make_media_link_replacer = (article) => {
  return (media) => {
    media.attrs
      .filter((a) => a.name == 'src')
      .forEach((a) => {
        const new_value = a.value.split('/').slice(-1)[0];
        a.value = `./${article.article_id}/${new_value}`;
      });
  };
};

///////////////////////////////////////////////
const user_id = 7308; // 自分のID

//api.collect_all_user_data(user_id).then((articles) => {
//  console.log(articles);
//}).catch((e) => {
//  console.log(e);
//});

//const article_id = 31509; // 動画あり
const article_id = 41299; // 画像あり

api.get_article(article_id).then((article) => {
  const doc = parse_contents(article.contents);
  const replace_media_link = make_media_link_replacer(article);
  for_each_media_nodes(doc, replace_media_link);

  console.dir(doc, {depth: null});
});
