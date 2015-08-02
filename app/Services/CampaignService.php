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
     * @var CampaignRepository
     */
    private $campaignRepository;

    public function __construct(CampaignRepository $repository)
    {
        $this->campaignRepository = $repository;
    }

    /**
     * Find all templates
     *
     * @param bool $paginate
     * @param int $perPage
     * @return mixed
     */
    public function findAllCampaigns($paginate = false, $perPage = 10)
    {
        if ($paginate) {
            return $this->campaignRepository->paginate($perPage);
        }

        return $this->campaignRepository->all();
    }

    /**
     * Find a template by id
     *
     * @param $id
     * @return mixed|null
     */
    public function findCampaign($id)
    {
        try {
            return $this->campaignRepository->find($id);
        } catch (ModelNotFoundException $e) {
            return null;
        }
    }

    /**
     * Create campaign
     *
     * @param array $data
     * @return mixed|null
     */
    public function createCampaign(array $data)
    {
        try {
            return $this->campaignRepository->create($data);
        } catch (Exception $e) {
            Log::error($e->getMessage());

            return null;
        }
    }

    /**
     * Update campaign by id
     *
     * @param array $data
     * @param $id
     * @return mixed|null
     */
    public function updateCampaign(array $data, $id)
    {
        try {
            return $this->campaignRepository->update($data, $id);
        } catch (Exception $e) {
            Log::error($e->getMessage());

            return null;
        }
    }

    /**
     * Delete a campaign by its id
     *
     * @param $campaignId
     * @return bool|int
     */
    public function deleteCampaign($campaignId)
    {
        try {
            return $this->campaignRepository->delete($campaignId);
        } catch (ModelNotFoundException $e) {

            return false;
        }
    }
}