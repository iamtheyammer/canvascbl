const { override, fixBabelImports } = require("customize-cra");

module.exports = override(
  fixBabelImports("import-antd", {
    libraryName: "antd",
    libraryDirectory: "es",
    style: "css"
  }),
  fixBabelImports("import-mobile", {
    libraryName: "antd-mobile",
    style: "css"
  }),
  fixBabelImports("import-lodash", {
    libraryName: "lodash",
    libraryDirectory: "",
    camel2DashComponentName: false // default: true
  })
);
