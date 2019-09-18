const glob = require("glob");
const markdownIt = require("markdown-it");
const meta = require("markdown-it-meta");
const ascii = require("./markdown-it-ascii.js");
const fs = require("fs");
const _ = require("lodash");
var path = require("path");

const sidebar = (directory, array) => {
  return array.map(i => {
    const children = _.sortBy(
      glob
        .sync(`./${directory}/${i[1]}/*.md`)
        .map(path => {
          const md = new markdownIt();
          const file = fs.readFileSync(path, "utf8");
          md.use(meta);
          md.render(file);
          const order = md.meta.order;
          return { path, order };
        })
        .filter(f => f.order !== false),
      ["order", "path"]
    ).map(f => f.path);
    return {
      title: i[0],
      children
    };
  });
};

module.exports = {
  title: "Cosmos SDK",
  base: process.env.VUEPRESS_BASE || "/",
  plugins: [
    [
      "@vuepress/search",
      {
        searchMaxSuggestions: 10
      }
    ]
  ],
  markdown: {
    anchor: {
      permalinkSymbol: ""
    }
    // extendMarkdown: md => {
    //   md.use(ascii);
    // }
  },
  locales: {
    "/": {
      lang: "en-US"
    },
    "/ru/": {
      lang: "ru"
    },
    "/kr/": {
      lang: "kr"
    },
    "/cn/": {
      lang: "cn"
    }
  },
  themeConfig: {
    repo: "cosmos/cosmos-sdk",
    docsDir: "docs",
    editLinks: true,
    locales: {
      "/": {
        label: "English",
        sidebar: sidebar("", [
          ["Intro", "intro"],
          ["Basics", "basics"],
          ["SDK Core", "core"],
          ["About Modules", "modules"],
          ["Interfaces", "interfaces"]
        ])
      },
      "/ru/": {
        label: "Русский",
        sidebar: sidebar("ru", [
          ["Введение", "intro"],
          ["Основы", "basics"],
          ["SDK Core", "core"],
          ["Модули", "modules"],
          ["Интерфейсы", "interfaces"]
        ])
      },
      "/kr/": {
        label: "한국어",
        sidebar: sidebar("kr", [
          ["소개", "intro"],
          ["기초", "basics"],
          ["SDK Core", "core"],
          ["모듈들", "modules"],
          ["인터페이스", "interfaces"]
        ])
      },
      "/cn/": {
        label: "中文",
        sidebar: sidebar("cn", [
          ["介绍", "intro"],
          ["基本", "basics"],
          ["SDK Core", "core"],
          ["模块", "modules"],
          ["接口", "interfaces"]
        ])
      }
    }
  }
};
