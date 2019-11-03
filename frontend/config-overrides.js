const { override, fixBabelImports } = require("customize-cra");

module.exports = override(
  fixBabelImports("import", {
    libraryName: "antd",
    libraryDirectory: "es",
    style: "css"
  }),
  fixBabelImports("import-mobile", {
    libraryName: "antd-mobile",
    style: "css"
  })
);
