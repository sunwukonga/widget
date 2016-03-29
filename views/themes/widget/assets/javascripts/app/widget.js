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

  var NAMESPACE = 'qor.widget';
  var EVENT_ENABLE = 'enable.' + NAMESPACE;
  var EVENT_DISABLE = 'disable.' + NAMESPACE;
  var EVENT_CLICK = 'click.' + NAMESPACE;
  var EDIT_WIDGET_BUTTON = '.js-widget-button, .qor-slideout__close';

  function QorWidget(element, options) {
    this.$element = $(element);
    this.options = $.extend({}, QorWidget.DEFAULTS, $.isPlainObject(options) && options);
    this.init();
  }

  QorWidget.prototype = {
    constructor: QorWidget,

    init: function () {
      var $this = this.$element;
      this.bind();
      this.initStatus();
    },

    bind: function () {
      this.$element.on(EVENT_CLICK, $.proxy(this.click, this));
    },

    initStatus : function () {
      $("body").append('<iframe id="qor-widget-iframe" src="/admin/widget_instances/frontend-edit"></iframe>');
    },

    click: function (e) {
      var $target = $(e.target);
      e.stopPropagation();

      if ($target.is(EDIT_WIDGET_BUTTON)){
        $("#qor-widget-iframe").contents().find(".js-widget-edit-link").attr("data-url", $target.data("url"));
        $("#qor-widget-iframe").addClass("show");
      }
    }
  };

  QorWidget.plugin = function (options) {
    return this.each(function () {
      var $this = $(this);
      var data = $this.data(NAMESPACE);
      var fn;

      if (!data) {

        if (/destroy/.test(options)) {
          return;
        }

        $this.data(NAMESPACE, (data = new QorWidget(this, options)));
      }

      if (typeof options === 'string' && $.isFunction(fn = data[options])) {
        fn.apply(data);
      }
    });
  };

  QorWidget.isScrollToBottom = function (element) {
    return element.clientHeight + element.scrollTop === element.scrollHeight;
  };

  $(function () {
    $("body").attr("data-toggle", "qor.widgets");
    var btnCss = 'style="background-color: rgb(33,150,243); color: white; border: none; box-shadow: 0 2px 2px 0 rgba(0,0,0,.14),0 3px 1px -2px rgba(0,0,0,.2),0 1px 5px 0 rgba(0,0,0,.12); padding: 6px 16px; display: inline-block; font-family: "Roboto","Helvetica","Arial",sans-serif; font-size: 14px; font-weight: 500;"';
    $(".qor-widget").each(function (i, e) {
      var $wrap = $(e).find("*").eq(0);
      $wrap.css("position", "relative").addClass("qor-widget").attr("data-url", $(e).data("url")).unwrap();
      $wrap.append('<div style="position: absolute; top: 8px; right: 8px;" class="qor-widget-embed-wrapper"><button ' + btnCss + ' data-url=\"' + $(e).data("url") + '\" class="js-widget-button">Edit Widget</button></div>');
    });
    window.closeWidgetEditBox = function () {
      $("#qor-widget-iframe").removeClass("show");
      $("#qor-widget-iframe")[0].contentWindow.location.reload();
    };

    var selector = '[data-toggle="qor.widgets"]';
    $(document).
      on(EVENT_DISABLE, function (e) {
        QorWidget.plugin.call($(selector, e.target), 'destroy');
      }).
      on(EVENT_ENABLE, function (e) {
        QorWidget.plugin.call($(selector, e.target));
      }).
      triggerHandler(EVENT_ENABLE);
  });

  return QorWidget;
});
