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
  var EVENT_CHANGE = 'change.' + NAMESPACE;
  var TARGET_WIDGET = '[name="QorResource.Widgets"]';

  function QorWidget(element, options) {
    this.$element = $(element);
    this.options = $.extend({}, QorWidget.DEFAULTS, $.isPlainObject(options) && options);
    this.init();
  }

  QorWidget.prototype = {
    constructor: QorWidget,

    init: function () {
      this.bind();
      this.addWidgetSlideout();
      this.initSelect();
    },

    bind: function () {
      this.$element.on(EVENT_CHANGE, 'select', this.change);
    },

    unbind: function () {
      this.$element.off(EVENT_CHANGE, 'select', this.change);
    },

    initSelect: function () {
      $('select').closest('.qor-form-section').hide();
      $('select').each(function () {
        // TODO:
        if ($(this).find('option').size() > 2) {
          $(this).closest('.qor-form-section').show();
        }
      });
    },

    addWidgetSlideout: function () {
      var $select = $(TARGET_WIDGET);

      $select.find('option').each(function () {
        var $this = $(this);
        var val = $this.val();

        var url = $select.closest('form').prop('action') + '?widget_type=' + val;
        if (val) {
          $select.after('<a href="#" style="display: none;" class="qor-widget-' + val + '" data-url="' + url + '">' + val + '</a>');
        }
      });
    },

    change: function (e) {
      var $target = $(e.target);
      var widgetValue = $target.val();
      var location = window.location;
      var isInSlideout = $target.closest('.qor-slideout').size();

      if (!$target.is(TARGET_WIDGET)) {
        return;
      }

      if (isInSlideout){
        var clickClass = '.qor-widget-' + widgetValue;
        $(clickClass).click();
      } else {

        if (widgetValue) {
          location.search = 'widget_type=' + widgetValue;
        }

      }
    },

    destroy: function () {
      this.unbind();
    }
  };

  QorWidget.DEFAULTS = {};

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

  $(function () {
    var selector = '[data-toggle="qor.widget"]';

    $(document)
      .on(EVENT_DISABLE, function (e) {
        QorWidget.plugin.call($(selector, e.target), 'destroy');
      })
      .on(EVENT_ENABLE, function (e) {
        QorWidget.plugin.call($(selector, e.target));
      })
      .triggerHandler(EVENT_ENABLE);
  });

  return QorWidget;

});
