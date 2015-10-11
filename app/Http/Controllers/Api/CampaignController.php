<?php

namespace newsletters\Http\Controllers\Api;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use newsletters\Http\Controllers\Controller;
use newsletters\Http\Requests\SendCampaignRequest;
use newsletters\Http\Requests\StoreCampaignRequest;
use newsletters\Http\Requests\TestSendRequest;
use newsletters\Jobs\SendCampaign;
use newsletters\Services\CampaignService;
use newsletters\Services\ListsService;
use newsletters\Services\UserService;

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
     * @return \Illuminate\Http\JsonResponse
     */
    public function index(Request $request)
    {
        $campaigns = $this->service->findAllCampaigns($request->has('paginate'), 10);

        return response()->json($campaigns, 200);
    }

    /**
     * Store a newly created resource in storage.
     *
     * @param StoreCampaignRequest $request
     * @return \Illuminate\Http\JsonResponse
     */
    public function store(StoreCampaignRequest $request)
    {
        $campaign = $this->service->createCampaign($request->all());
        if (isset($campaign)) {
            return response()->json(['campaign' => $campaign->id], 200);
        }

        return response()->json(['message' => ['The specified resource could not be created.']], 412);
    }

    /**
     * Display the specified resource.
     *
     * @param  int $id
     * @return \Illuminate\Http\JsonResponse
     */
    public function show($id)
    {
        $campaign = $this->service->findCampaign($id);

        if (isset($campaign)) {
            return response()->json($campaign, 200);
        }

        return response()->json(['message' => ['The specified resource does not exist.']], 404);
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
        $campaign = $this->service->updateCampaign($request->all(), $id);
        if (isset($campaign)) {
            return response()->json(['campaign' => $campaign->id], 200);
        }

        return response()->json(['message' => ['The specified resource could not be updated.']], 412);
    }

    /**
     * Remove the specified resource from storage.
     *
     * @param  int $id
     * @return Response
     */
    public function destroy($id)
    {
        if ($this->service->deleteCampaign($id)) {
            return response()->json(['message' => ['The specified resource has been deleted.']], 200);
        }

        return response()->json(['message' => ['The specified resource could not be deleted.']], 422);
    }

    /**
     * Send campaign
     *
     * @param SendCampaignRequest $request
     * @param UserService $userService
     * @param ListsService $listsService
     * @return \Illuminate\Http\JsonResponse
     */
    public function send(SendCampaignRequest $request, UserService $userService, ListsService $listsService)
    {
        $campaign = $this->service->findCampaign($request->input('id'));
        $subscribers = $listsService->findAllSubscribersByListIds($request->input('lists'));

        $user = Auth::user();
        $userService->setSesConfig($user->aws_key, $user->aws_secret, $user->aws_region);

        $this->dispatch(new SendCampaign($campaign, $subscribers));

        return response()->json(['message' => ['The campaign has been started.']], 200);
    }

    public function testSend(TestSendRequest $request)
    {
        //TODO Test send the campaign to the emails specified in the request
    }
}
