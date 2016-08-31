(function (factory) {
  if (typeof define === 'function' && define.amd) {
    // AMD. Register as anonymous module.
    define(['jquery'], factory);
  } else if (typeof exports === 'object') {
    // Node / CommonJS
    factory(require('jquery'));
  } else {
    // Browser globals.
    factory(jQuery);
  }
})(function ($) {

  'use strict';

  var $body = $("body");
  var NAMESPACE = 'qor.widget.inlineEdit';
  var EVENT_ENABLE = 'enable.' + NAMESPACE;
  var EVENT_DISABLE = 'disable.' + NAMESPACE;
  var EVENT_CLICK = 'click.' + NAMESPACE;
  var EDIT_WIDGET_BUTTON = '.qor-widget-button';
  var ID_WIDGET = '#qor-widget-iframe';
  var INLINE_EDIT_URL;

  function QorWidgetInlineEdit(element, options) {
    this.$element = $(element);
    this.options = $.extend({}, QorWidgetInlineEdit.DEFAULTS, $.isPlainObject(options) && options);
    this.init();
  }

  QorWidgetInlineEdit.prototype = {
    constructor: QorWidgetInlineEdit,

    init: function () {
      this.bind();
      this.initStatus();
    },

    bind: function () {
      this.$element.on(EVENT_CLICK, EDIT_WIDGET_BUTTON, this.click);
    },

    initStatus : function () {
      $body.append('<iframe id="qor-widget-iframe" src="' + INLINE_EDIT_URL + '"></iframe>');
    },

    click: function () {
      var $this = $(this);

      $(ID_WIDGET).addClass("show").focus();
      document.getElementById('qor-widget-iframe').contentWindow.$(".js-widget-edit-link").attr("data-url", $this.data("url")).click();
      $body.addClass("open-widget-editor");

      return false;
    }
  };

  QorWidgetInlineEdit.plugin = function (options) {
    return this.each(function () {
      var $this = $(this);
      var data = $this.data(NAMESPACE);
      var fn;

      if (!data) {

        if (/destroy/.test(options)) {
          return;
        }

        $this.data(NAMESPACE, (data = new QorWidgetInlineEdit(this, options)));
      }

      if (typeof options === 'string' && $.isFunction(fn = data[options])) {
        fn.apply(data);
      }
    });
  };


  $(function () {
    $body.attr("data-toggle", "qor.widgets");
    $(".qor-widget").each(function () {
      var $this = $(this);
      var $wrap = $this.children().eq(0);
      INLINE_EDIT_URL = $this.data("widget-inline-edit-url");
      $wrap.css("position", "relative").addClass("qor-widget").unwrap();
      $wrap.append('<div class="qor-widget-embed-wrapper"><button data-url=\"' + $this.data("url") + '\" class="qor-widget-button">Edit</button></div>');
    });
    
    window.closeWidgetEditBox = function () {
      $(ID_WIDGET).removeClass("show");
      $body.removeClass("open-widget-editor");
      $(ID_WIDGET)[0].contentWindow.location.reload();
    };

    var selector = '[data-toggle="qor.widgets"]';
    $(document).
      on(EVENT_DISABLE, function (e) {
        QorWidgetInlineEdit.plugin.call($(selector, e.target), 'destroy');
      }).
      on(EVENT_ENABLE, function (e) {
        QorWidgetInlineEdit.plugin.call($(selector, e.target));
      }).
      triggerHandler(EVENT_ENABLE);
  });


  return QorWidgetInlineEdit;
});
