@inject('templates', 'newsletters\Services\TemplateService')

{!! $templates->renderTemplate($template_id, $name, $custom_fields) !!}