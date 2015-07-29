<?php

namespace newsletters\Http\Controllers\Api;

use Illuminate\Http\Request;
use newsletters\Http\Controllers\Controller;
use newsletters\Http\Requests;
use newsletters\Http\Requests\StoreCampaignRequest;
use newsletters\Repositories\CampaignRepository;
use newsletters\Services\CampaignService;

class CampaignController extends Controller
{
    /**
     * @var CampaignService
     */
    private $service;

    public function __construct(CampaignService $service)
    {
        $this->middleware('auth.basic');

        $this->service = $service;
    }

    /**
     * Display a listing of the resource.
     *
     * @param Request $request
     * @param CampaignRepository $repository
     * @return \Illuminate\Http\JsonResponse
     */
    public function index(Request $request, CampaignRepository $repository)
    {
        $campaigns = $this->service->findAllCampaigns($request->has('paginate'), 10, $repository);

        return response()->json($campaigns, 200);
    }

    /**
     * Store a newly created resource in storage.
     *
     * @param StoreCampaignRequest $request
     * @param CampaignRepository $repository
     * @return \Illuminate\Http\JsonResponse
     */
    public function store(StoreCampaignRequest $request, CampaignRepository $repository)
    {
        $campaign = $this->service->createCampaign($request->all(), $repository);
        if (isset($campaign)) {
            return response()->json(['status' => 200, 'campaign' => $campaign], 200);
        }

        return response()->json(['status' => 412, 'campaign' => ['The specified resource could not be created.']],
            412);
    }

    /**
     * Display the specified resource.
     *
     * @param  int $id
     * @param CampaignRepository $repository
     * @return \Illuminate\Http\JsonResponse
     */
    public function show($id, CampaignRepository $repository)
    {
        $campaign = $this->service->findCampaign($id, $repository);

        if (isset($campaign)) {
            return response()->json($campaign, 200);
        }

        return response()->json(['status' => 404, 'message' => 'The specified resource does not exist.'], 404);
    }

    /**
     * Update the specified resource in storage.
     *
     * @param  Request $request
     * @param  int $id
     * @return Response
     */
    public function update(Request $request, $id)
    {
        //
    }

    /**
     * Remove the specified resource from storage.
     *
     * @param  int $id
     * @param CampaignRepository $repository
     * @return Response
     */
    public function destroy($id, CampaignRepository $repository)
    {
        if ($this->service->deleteCampaign($id, $repository)) {
            return response()->json(['status' => 200, 'message' => 'The specified resource has been deleted.'],
                200);
        }

        return response()->json(['status' => 422, 'campaign' => ['The specified resource could not be deleted.']],
            422);

    }
}
