<?php
/**
 * Created by PhpStorm.
 * User: filip
 * Date: 27.7.15
 * Time: 21:36
 */

namespace newsletters\Services;

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
     * Find all campaigns
     *
     * @return mixed
     */
    public function findAllCampaigns()
    {
        return $this->campaignRepository->all();
    }

    /**
     * Find all campaigns in a paginated list
     *
     * @param $perPage
     * @return mixed
     */
    public function findAllCampaignsPaginated($perPage = 10)
    {
        return $this->campaignRepository->paginate($perPage);
    }

    /**
     * Find a campaign by id
     *
     * @param $id
     * @return mixed|null
     */
    public function findCampaign($id, array $with = [])
    {
        return $this->campaignRepository->with($with)->find($id);
    }

    /**
     * Create campaign
     *
     * @param array $data
     * @return mixed
     */
    public function createCampaign(array $data)
    {
        return $this->campaignRepository->create($data);
    }

    /**
     * Update campaign by id
     *
     * @param array $data
     * @param $id
     * @return mixed
     */
    public function updateCampaign(array $data, $id)
    {
        return $this->campaignRepository->update($data, $id);
    }

    /**
     * Delete a campaign by its id
     *
     * @param $id
     * @return bool|int
     */
    public function deleteCampaign($id)
    {
        return $this->campaignRepository->delete($id);
    }
}
