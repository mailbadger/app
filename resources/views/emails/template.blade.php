@inject('templates', 'newsletters\Services\TemplateService')

{!! $templates->renderTemplate($template_id, $name, $email, $custom_fields) !!}
