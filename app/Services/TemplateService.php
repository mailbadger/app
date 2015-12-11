<?php

namespace newsletters\Services;

use Illuminate\Database\Eloquent\ModelNotFoundException;
use Illuminate\Support\Facades\Log;
use newsletters\Repositories\TemplateRepository;
use Sunra\PhpSimple\HtmlDomParser;

/**
 * Created by PhpStorm.
 * User: filip
 * Date: 27.7.15
 * Time: 21:30
 */
class TemplateService
{
    /**
     * @var TemplateRepository
     */
    private $templateRepository;

    public function __construct(TemplateRepository $repository)
    {
        $this->templateRepository = $repository;
    }

    /**
     * Find all templates
     *
     * @return mixed
     */
    public function findAllTemplates()
    {
        return $this->templateRepository->all(['id', 'name', 'created_at', 'updated_at']);
    }

    /**
     * Find all templates paginated
     *
     * @param int $perPage
     * @return mixed
     */
    public function findAllTemplatesPaginated($perPage = 10)
    {
        return $this->templateRepository->paginate($perPage, ['id', 'name', 'created_at', 'updated_at']); 
    }

    /**
     * Find a template by id
     *
     * @param $id
     * @return mixed|null
     */
    public function findTemplate($id)
    {
        return $this->templateRepository->find($id);
    }

    /**
     * Create template
     *
     * @param array $data
     * @return mixed|null
     */
    public function createTemplate(array $data)
    {
        return $this->templateRepository->create($data);
    }

    /**
     * Update template by id
     *
     * @param array $data
     * @param $id
     * @return mixed|null
     */
    public function updateTemplate(array $data, $id)
    {
        return $this->templateRepository->update($data, $id);
    }

    /**
     * Deletes a template if it's unused by a campaign
     *
     * @param $templateId
     * @return int
     */
    public function deleteUnusedTemplate($templateId)
    {
        try {
            return $this->templateRepository
                ->scopeQuery(function ($query) {
                    return $query->has('campaigns', '=', 0);
                })
                ->find($templateId)
                ->delete();
        } catch (ModelNotFoundException $e) {
            return 0;
        }
    }

    /**
     * Render an html template with the custom fields replaced
     *
     * @param $templateId
     * @param $subscriberName
     * @param $subscriberEmail
     * @param $opensTrackerUrl
     * @param array $customTags
     * @return string
     */
    public function renderTemplate($templateId, $subscriberName, $subscriberEmail, $opensTrackerUrl, array $customTags = [])
    { 
        $template = $this->templateRepository->find($templateId);  

        $dom = HtmlDomParser::str_get_html($template->content);

        $dom = $this->replaceTagsInTemplate($dom, $customTags); 
        
        $dom = $this->appendImg($dom, $opensTrackerUrl);

        $html = $dom->outertext;

        $dom->clear();
        unset($dom);

        return $html;
    }

    /**
     * Replace the personalization tags that are put in the template content
     * eg tags: *|Name|*, *|Email|* etc..
     *
     * @param $content
     * @param array $tags
     * @return Dom 
     */
    private function replaceTagsInTemplate($dom, array $tags)
    {
        foreach($tags as $key => $val) {
            $dom->outertext = preg_replace($key, $val, $dom->outertext);
        }

        return $dom;
    }

    /**
     * Append image tag used for tracking the email opens
     *
     * @param $content
     * @param $url
     * @return string
     */
    private function appendImg($dom, $url)
    {
        $img = '<img src="'.$url.'"/>'.PHP_EOL;
 
        //append the image inside the body tag if the dom has it
        $body = $dom->find('html body', 0);
 
        if(!empty($body)) {
            $body->innertext .= $img;
        } else {
            $dom->outertext .= $img;
        }
  
        return $dom;
    }
}
