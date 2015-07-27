<?php

namespace newsletters\Services;

use Illuminate\Database\Eloquent\ModelNotFoundException;
use newsletters\Repositories\TemplateRepository;

/**
 * Created by PhpStorm.
 * User: filip
 * Date: 27.7.15
 * Time: 21:30
 */
class TemplateService
{

    /**
     * Find all templates
     *
     * @param bool $paginate
     * @param int $perPage
     * @param TemplateRepository $repository
     * @return mixed
     */
    public function findAllTemplates($paginate = false, $perPage = 10, TemplateRepository $repository)
    {
        if ($paginate) {
            return $repository->paginate($perPage);
        }

        return $repository->all();
    }

    /**
     * Find a template by id
     *
     * @param $id
     * @param TemplateRepository $repository
     * @return mixed|null
     */
    public function findTemplate($id, TemplateRepository $repository)
    {
        try {
            return $repository->find($id);
        } catch (ModelNotFoundException $e) {
            return null;
        }
    }

    /**
     * Deletes a template if it's unused by a campaign
     *
     * @param $templateId
     * @param TemplateRepository $repository
     * @return bool|int
     */
    public function deleteUnusedTemplate($templateId, TemplateRepository $repository)
    {
        try {
            return $repository
                ->scopeQuery(function ($query) {
                    return $query->has('campaigns', '=', 0);
                })
                ->find($templateId)
                ->delete();
        } catch (ModelNotFoundException $e) {
            return false;
        }
    }
}