<?php

namespace newsletters\Services;

use Exception;
use Illuminate\Database\Eloquent\ModelNotFoundException;
use Illuminate\Support\Facades\Log;
use newsletters\Repositories\TemplateRepository;
use Symfony\Component\DomCrawler\Crawler;

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
     * @param bool $paginate
     * @param int $perPage
     * @return mixed
     */
    public function findAllTemplates($paginate = false, $perPage = 10)
    {
        $columns = ['id', 'name', 'created_at', 'updated_at'];

        return (!empty($paginate))
            ? $this->templateRepository->paginate($perPage, $columns)
            : $this->templateRepository->all($columns);
    }

    /**
     * Find a template by id
     *
     * @param $id
     * @return mixed|null
     */
    public function findTemplate($id)
    {
        try {
            return $this->templateRepository->find($id);
        } catch (ModelNotFoundException $e) {
            return null;
        }
    }

    /**
     * Create template
     *
     * @param array $data
     * @return mixed|null
     */
    public function createTemplate(array $data)
    {
        try {
            return $this->templateRepository->create($data);
        } catch (Exception $e) {
            Log::error($e->getMessage());

            return null;
        }
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
        try {
            return $this->templateRepository->update($data, $id);
        } catch (Exception $e) {
            Log::error($e->getMessage());

            return null;
        }
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
     * @param array $customFields
     * @return string
     */
    public function renderTemplate($templateId, $subscriberName, array $customFields)
    {
        $template = $this->templateRepository->find($templateId);
        $content = $template->content;
        $content = str_replace('*|Name|*', $subscriberName, $content);

        foreach ($customFields as $key => $val) {
            //TODO replace with preg_replace for other possible matches in the template such as *|Foo|* *|foo|* etc..
            $content = str_replace('*|' . $key . '|*', $val, $content);
        }

        $html = new Crawler($content);

        return $html->html();
    }
}
