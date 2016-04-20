if(!window.loadedWidgetAsset) {
  window.loadedWidgetAsset = true;
  var prefix = document.currentScript.getAttribute("data-prefix");
  document.write("<script src='" + prefix + "/assets/javascripts/vendors/jquery.min.js'></script><script src=\"" + prefix + "/assets/javascripts/widget_inline_edit.js?theme=widget\"></script><link type=\"text/css\" rel=\"stylesheet\" href=\"" + prefix + "/assets/stylesheets/widget_inline_edit.css?theme=widget\">");
}
