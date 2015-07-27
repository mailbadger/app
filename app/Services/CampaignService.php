<?php
/**
 * Created by PhpStorm.
 * User: filip
 * Date: 27.7.15
 * Time: 21:36
 */

namespace newsletters\Services;


use Exception;
use Illuminate\Database\Eloquent\ModelNotFoundException;
use Illuminate\Support\Facades\Log;
use newsletters\Repositories\CampaignRepository;

class CampaignService
{

    /**
     * Find all templates
     *
     * @param bool $paginate
     * @param int $perPage
     * @param CampaignRepository $repository
     * @return mixed
     */
    public function findAllCampaigns($paginate = false, $perPage = 10, CampaignRepository $repository)
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
     * @param CampaignRepository $repository
     * @return mixed|null
     */
    public function findCampaign($id, CampaignRepository $repository)
    {
        try {
            return $repository->find($id);
        } catch (ModelNotFoundException $e) {
            return null;
        }
    }

    /**
     * Create campaign
     *
     * @param array $data
     * @param CampaignRepository $repository
     * @return mixed|null
     */
    public function createCampaign(array $data, CampaignRepository $repository)
    {
        try {
            return $repository->create($data);
        } catch (Exception $e) {
            Log::error($e->getMessage());

            return null;
        }
    }
    /**
     * Delete a campaign by its id
     *
     * @param $campaignId
     * @param CampaignRepository $repository
     * @return bool|int
     */
    public function deleteCampaign($campaignId, CampaignRepository $repository)
    {
        try {
            return $repository->delete($campaignId);
        } catch (ModelNotFoundException $e) {

            return false;
        }
    }
}